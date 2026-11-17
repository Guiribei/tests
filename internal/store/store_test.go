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
	if err := db.SaveAnswer(ctx, "session-1", attempt.ID, 1, fixture.Document.AnswerKey["1"]); err != nil {
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
	if err := db.Finish(ctx, "session-1", attempt.ID, score); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAnswer(ctx, "session-1", attempt.ID, 2, "A"); !errors.Is(err, ErrFinished) {
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
