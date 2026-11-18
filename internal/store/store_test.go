package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/scoring"
)

func TestAttemptLifecycle(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	fixture, err := exam.LoadFile("../../option-one.json")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	attempt, err := db.StartOrResume(ctx, "attempt-1", "session-1", fixture)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ApplyAnswers(ctx, "session-1", attempt.ID, nil, map[int]string{1: fixture.Document.AnswerKey["1"]}, false); err != nil {
		t.Fatal(err)
	}
	answers, err := db.Answers(ctx, attempt.ID)
	if err != nil {
		t.Fatal(err)
	}
	score := scoring.Calculate(fixture, answers)
	if score.TotalCorrect != 1 {
		t.Fatalf("expected one correct answer, got %d", score.TotalCorrect)
	}
	if _, err := db.ApplyAnswers(ctx, "session-1", attempt.ID, nil, nil, true); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ApplyAnswers(ctx, "session-1", attempt.ID, nil, map[int]string{2: "A"}, false); !errors.Is(err, ErrFinished) {
		t.Fatalf("expected ErrFinished, got %v", err)
	}
	cohort, err := db.Cohort(ctx, fixture.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if len(cohort) != 1 || cohort[0].TotalCorrect != 1 {
		t.Fatalf("unexpected cohort: %#v", cohort)
	}
}

func TestSyncAndFinalRetry(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "offline.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	e, err := exam.LoadFile("../../option-three.json")
	if err != nil {
		t.Fatal(err)
	}
	a, err := db.StartOrResume(ctx, "attempt", "owner", e)
	if err != nil {
		t.Fatal(err)
	}

	var initial int64
	draft, err := db.ApplyAnswers(ctx, "owner", a.ID, &initial, map[int]string{1: "A", 2: "B"}, false)
	if err != nil || draft.Revision != 1 {
		t.Fatalf("draft: %+v, %v", draft, err)
	}
	stale, err := db.ApplyAnswers(ctx, "owner", a.ID, &initial, map[int]string{1: "D"}, true)
	if !errors.Is(err, ErrRevision) || stale.Answers[1] != "A" || stale.Finished {
		t.Fatalf("stale final: %+v, %v", stale, err)
	}
	if _, err := db.ApplyAnswers(ctx, "other", a.ID, &draft.Revision, map[int]string{}, true); !errors.Is(err, ErrNotFound) {
		t.Fatalf("foreign attempt: %v", err)
	}
	if _, err := db.ApplyAnswers(ctx, "owner", a.ID, &draft.Revision, map[int]string{2: "E", 181: "A"}, true); !errors.Is(err, ErrInvalidAnswer) {
		t.Fatalf("invalid submission: %v", err)
	}
	_, before, err := db.Packet(ctx, "owner", a.ID)
	if err != nil || before.Revision != draft.Revision || before.Answers[2] != "B" {
		t.Fatalf("partial submission: %+v, %v", before, err)
	}
	choices := map[int]string{1: e.Document.AnswerKey["1"], 81: e.Document.AnswerKey["81"]}
	final, err := db.ApplyAnswers(ctx, "owner", a.ID, &draft.Revision, choices, true)
	if err != nil || !final.Finished || final.Revision != 2 || len(final.Answers) != 2 {
		t.Fatalf("final: %+v, %v", final, err)
	}
	repeated, err := db.ApplyAnswers(ctx, "owner", a.ID, &draft.Revision, map[int]string{}, true)
	if err != nil || !repeated.Finished || repeated.Revision != final.Revision || repeated.Answers[81] != choices[81] {
		t.Fatalf("duplicate final: %+v, %v", repeated, err)
	}
	got, err := db.GetAttempt(ctx, "owner", a.ID)
	if err != nil || got.WeightedRaw != 3 || got.FinishedAt == nil {
		t.Fatalf("score=%d err=%v", got.WeightedRaw, err)
	}
}

func TestUpgradeActiveAttempt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	e, err := exam.LoadFile("../../option-two.json")
	if err != nil {
		t.Fatal(err)
	}
	a, err := db.StartOrResume(ctx, "existing", "owner", e)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ApplyAnswers(ctx, "owner", a.ID, nil, map[int]string{23: "B"}, false); err != nil {
		t.Fatal(err)
	}
	// Reproduce the deployed schema, which predates the revision column.
	if _, err := db.db.Exec("ALTER TABLE attempts DROP COLUMN revision; PRAGMA user_version=0"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	snapshot, sheet, err := db.Packet(ctx, "owner", a.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Document.Metadata.Title != a.ExamTitle || sheet.Revision != 0 || sheet.Finished || sheet.Answers[23] != "B" {
		t.Fatal("tentativa antiga alterada na migração")
	}
}
