package scoring

import (
	"math"
	"testing"
)

func TestStandardize(t *testing.T) {
	cohort := []Score{
		{GeneralCorrect: 40, SpecificCorrect: 50},
		{GeneralCorrect: 50, SpecificCorrect: 60},
		{GeneralCorrect: 60, SpecificCorrect: 70},
	}
	result := Standardize(cohort[1], cohort)
	if !result.Available {
		t.Fatal("expected standardized score")
	}
	if math.Abs(result.GeneralStandard-50) > 0.001 || math.Abs(result.SpecificStandard-50) > 0.001 {
		t.Fatalf("expected mean candidate to score 50/50, got %.2f/%.2f", result.GeneralStandard, result.SpecificStandard)
	}
	if math.Abs(result.Final-150) > 0.001 || !result.Qualified {
		t.Fatalf("expected final 150 and qualified, got %.2f", result.Final)
	}
}

func TestStandardizeNeedsCohortVariance(t *testing.T) {
	result := Standardize(Score{GeneralCorrect: 50, SpecificCorrect: 60}, []Score{{GeneralCorrect: 50, SpecificCorrect: 60}})
	if result.Available {
		t.Fatal("single candidate must not produce a standardized score")
	}
}
