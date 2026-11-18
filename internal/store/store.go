package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/scoring"
	_ "modernc.org/sqlite"
)

var (
	ErrNotFound      = errors.New("tentativa não encontrada")
	ErrRevision      = errors.New("revisão desatualizada")
	ErrInvalidAnswer = errors.New("resposta inválida")
	ErrFinished      = errors.New("tentativa já finalizada")
)

type Store struct{ db *sql.DB }

type Attempt struct {
	ID              string
	SessionID       string
	ExamSlug        string
	ExamTitle       string
	ExamJSON        []byte
	Revision        int64
	Status          string
	StartedAt       time.Time
	UpdatedAt       time.Time
	FinishedAt      *time.Time
	GeneralCorrect  int
	SpecificCorrect int
	WeightedRaw     int
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	store := &Store{db: db}
	if err := store.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS attempts (
  id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL,
  exam_slug TEXT NOT NULL,
  exam_title TEXT NOT NULL,
  exam_json BLOB NOT NULL,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'finished')),
  started_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  finished_at TEXT,
  general_correct INTEGER NOT NULL DEFAULT 0,
  specific_correct INTEGER NOT NULL DEFAULT 0,
  weighted_raw INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS one_active_attempt
  ON attempts(session_id, exam_slug) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS attempts_exam_status ON attempts(exam_slug, status);
CREATE TABLE IF NOT EXISTS answers (
  attempt_id TEXT NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
  question_number INTEGER NOT NULL,
  choice TEXT NOT NULL CHECK (choice IN ('A','B','C','D','E')),
  updated_at TEXT NOT NULL,
  PRIMARY KEY (attempt_id, question_number)
);`
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, schema); err != nil {
		return err
	}
	var version int
	if err := tx.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version == 0 {
		if _, err := tx.ExecContext(ctx, "ALTER TABLE attempts ADD COLUMN revision INTEGER NOT NULL DEFAULT 0; PRAGMA user_version=1"); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) StartOrResume(ctx context.Context, id, sessionID string, e *exam.Exam) (Attempt, error) {
	existing, err := s.activeAttempt(ctx, sessionID, e.Slug)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return Attempt{}, err
	}
	snapshot, err := e.Snapshot()
	if err != nil {
		return Attempt{}, err
	}
	now := time.Now().UTC()
	_, err = s.db.ExecContext(ctx, `INSERT INTO attempts
      (id, session_id, exam_slug, exam_title, exam_json, started_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?)`, id, sessionID, e.Slug, e.Document.Metadata.Title, snapshot, formatTime(now), formatTime(now))
	if err != nil {
		return Attempt{}, err
	}
	return s.GetAttempt(ctx, sessionID, id)
}

func (s *Store) activeAttempt(ctx context.Context, sessionID, slug string) (Attempt, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, session_id, exam_slug, exam_title, exam_json, status, revision,
      started_at, updated_at, finished_at, general_correct, specific_correct, weighted_raw
      FROM attempts WHERE session_id = ? AND exam_slug = ? AND status = 'active'`, sessionID, slug)
	return scanAttempt(row)
}

func (s *Store) GetAttempt(ctx context.Context, sessionID, id string) (Attempt, error) {
	return getAttempt(ctx, s.db, sessionID, id)
}

type reader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func getAttempt(ctx context.Context, db reader, sessionID, id string) (Attempt, error) {
	row := db.QueryRowContext(ctx, `SELECT id, session_id, exam_slug, exam_title, exam_json, status, revision,
      started_at, updated_at, finished_at, general_correct, specific_correct, weighted_raw
      FROM attempts WHERE id = ? AND session_id = ?`, id, sessionID)
	return scanAttempt(row)
}

type scanner interface{ Scan(...any) error }

func scanAttempt(row scanner) (Attempt, error) {
	var attempt Attempt
	var started, updated string
	var finished sql.NullString
	err := row.Scan(&attempt.ID, &attempt.SessionID, &attempt.ExamSlug, &attempt.ExamTitle, &attempt.ExamJSON,
		&attempt.Status, &attempt.Revision, &started, &updated, &finished, &attempt.GeneralCorrect, &attempt.SpecificCorrect, &attempt.WeightedRaw)
	if errors.Is(err, sql.ErrNoRows) {
		return Attempt{}, ErrNotFound
	}
	if err != nil {
		return Attempt{}, err
	}
	attempt.StartedAt, err = time.Parse(time.RFC3339Nano, started)
	if err != nil {
		return Attempt{}, err
	}
	attempt.UpdatedAt, err = time.Parse(time.RFC3339Nano, updated)
	if err != nil {
		return Attempt{}, err
	}
	if finished.Valid {
		value, err := time.Parse(time.RFC3339Nano, finished.String)
		if err != nil {
			return Attempt{}, err
		}
		attempt.FinishedAt = &value
	}
	return attempt, nil
}

func (s *Store) Answers(ctx context.Context, attemptID string) (map[int]string, error) {
	return readAnswers(ctx, s.db, attemptID)
}

func readAnswers(ctx context.Context, db reader, attemptID string) (map[int]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT question_number, choice FROM answers WHERE attempt_id = ?`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	answers := make(map[int]string)
	for rows.Next() {
		var number int
		var choice string
		if err := rows.Scan(&number, &choice); err != nil {
			return nil, err
		}
		answers[number] = choice
	}
	return answers, rows.Err()
}

func (s *Store) Cohort(ctx context.Context, examSlug string) ([]scoring.Score, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT general_correct, specific_correct, weighted_raw
      FROM attempts WHERE exam_slug = ? AND status = 'finished'`, examSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scores []scoring.Score
	for rows.Next() {
		var score scoring.Score
		if err := rows.Scan(&score.GeneralCorrect, &score.SpecificCorrect, &score.WeightedRaw); err != nil {
			return nil, err
		}
		score.TotalCorrect = score.GeneralCorrect + score.SpecificCorrect
		score.WeightedMaximum = 280
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

func formatTime(value time.Time) string { return value.Format(time.RFC3339Nano) }
