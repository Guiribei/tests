package scoring

import (
	"fmt"
	"math"

	"github.com/Guiribei/tests/internal/exam"
)

type Score struct {
	GeneralCorrect  int
	SpecificCorrect int
	TotalCorrect    int
	WeightedRaw     int
	WeightedMaximum int
}

type CohortScore struct {
	Available         bool
	Candidates        int
	GeneralMean       float64
	GeneralDeviation  float64
	SpecificMean      float64
	SpecificDeviation float64
	GeneralStandard   float64
	SpecificStandard  float64
	Final             float64
	Qualified         bool
}

func Calculate(e *exam.Exam, answers map[int]string) Score {
	score := Score{WeightedMaximum: 280}
	for _, question := range e.Questions {
		correct := e.Document.AnswerKey[fmt.Sprint(question.Number)]
		if answers[question.Number] != correct {
			continue
		}
		score.TotalCorrect++
		score.WeightedRaw += question.Weight
		if question.Weight == 1 {
			score.GeneralCorrect++
		} else {
			score.SpecificCorrect++
		}
	}
	return score
}

func Standardize(candidate Score, cohort []Score) CohortScore {
	result := CohortScore{Candidates: len(cohort)}
	if len(cohort) < 2 {
		return result
	}
	general := make([]float64, len(cohort))
	specific := make([]float64, len(cohort))
	for i, score := range cohort {
		general[i] = float64(score.GeneralCorrect)
		specific[i] = float64(score.SpecificCorrect)
	}
	result.GeneralMean, result.GeneralDeviation = populationStats(general)
	result.SpecificMean, result.SpecificDeviation = populationStats(specific)
	if result.GeneralDeviation == 0 || result.SpecificDeviation == 0 {
		return result
	}
	result.GeneralStandard = ((float64(candidate.GeneralCorrect)-result.GeneralMean)/result.GeneralDeviation)*10 + 50
	result.SpecificStandard = ((float64(candidate.SpecificCorrect)-result.SpecificMean)/result.SpecificDeviation)*10 + 50
	result.Final = result.GeneralStandard + result.SpecificStandard*2
	result.Qualified = result.Final >= 150
	result.Available = true
	return result
}

func populationStats(values []float64) (float64, float64) {
	var total float64
	for _, value := range values {
		total += value
	}
	mean := total / float64(len(values))
	var squared float64
	for _, value := range values {
		difference := value - mean
		squared += difference * difference
	}
	return mean, math.Sqrt(squared / float64(len(values)))
}
