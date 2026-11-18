package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/scoring"
)

type Sheet struct {
	Revision int64          `json:"revision"`
	Answers  map[int]string `json:"answers"`
	Finished bool           `json:"finished"`
}

func readSheet(ctx context.Context, tx *sql.Tx, sessionID, id string) (Attempt, Sheet, error) {
	a, err := getAttempt(ctx, tx, sessionID, id)
	if err != nil {
		return Attempt{}, Sheet{}, err
	}
	answers, err := readAnswers(ctx, tx, id)
	return a, Sheet{Revision: a.Revision, Finished: a.Status == "finished", Answers: answers}, err
}

// Packet reads the immutable exam and a consistent revision of its answers.
func (s *Store) Packet(ctx context.Context, sessionID, id string) (*exam.Exam, Sheet, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, Sheet{}, err
	}
	defer tx.Rollback()
	a, sheet, err := readSheet(ctx, tx, sessionID, id)
	if err != nil {
		return nil, Sheet{}, err
	}
	e, err := exam.FromSnapshot(a.ExamJSON, a.ExamSlug)
	return e, sheet, err
}

// A nil revision is used by the online HTML forms. The offline client sends
// changes for a draft and the entire answer sheet for a final submission.
func (s *Store) ApplyAnswers(ctx context.Context, sessionID, id string, revision *int64, choices map[int]string, finish bool) (Sheet, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Sheet{}, err
	}
	defer tx.Rollback()
	a, sheet, err := readSheet(ctx, tx, sessionID, id)
	if err != nil {
		return Sheet{}, err
	}
	if sheet.Finished {
		if revision == nil && !finish {
			return sheet, ErrFinished
		}
		return sheet, nil
	}
	if revision != nil && *revision != sheet.Revision {
		return sheet, ErrRevision
	}
	e, err := exam.FromSnapshot(a.ExamJSON, a.ExamSlug)
	if err != nil {
		return Sheet{}, err
	}
	for n, choice := range choices {
		if _, ok := e.Question(n); !ok || len(choice) != 1 || choice < "A" || choice > "E" {
			return Sheet{}, ErrInvalidAnswer
		}
	}
	if finish && revision != nil {
		if _, err := tx.ExecContext(ctx, "DELETE FROM answers WHERE attempt_id = ?", id); err != nil {
			return Sheet{}, err
		}
		sheet.Answers = make(map[int]string)
	}
	now := formatTime(time.Now().UTC())
	for n, choice := range choices {
		if _, err := tx.ExecContext(ctx, `INSERT INTO answers (attempt_id, question_number, choice, updated_at)
			VALUES (?, ?, ?, ?) ON CONFLICT(attempt_id, question_number) DO UPDATE SET
			choice=excluded.choice, updated_at=excluded.updated_at`, id, n, choice, now); err != nil {
			return Sheet{}, err
		}
		sheet.Answers[n] = choice
	}
	if finish {
		score := scoring.Calculate(e, sheet.Answers)
		_, err = tx.ExecContext(ctx, `UPDATE attempts SET status='finished', finished_at=?, updated_at=?,
			general_correct=?, specific_correct=?, weighted_raw=?, revision=revision+1 WHERE id=?`,
			now, now, score.GeneralCorrect, score.SpecificCorrect, score.WeightedRaw, id)
	} else {
		_, err = tx.ExecContext(ctx, "UPDATE attempts SET updated_at=?, revision=revision+1 WHERE id=?", now, id)
	}
	if err != nil {
		return Sheet{}, err
	}
	sheet.Revision++
	sheet.Finished = finish
	return sheet, tx.Commit()
}
