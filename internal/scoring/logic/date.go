package logic

import (
	"math"
	"techsupport/internal/models"
	"techsupport/internal/scoring"
	"time"
)

func CalculateDeltaOfCreationAge(userClaim time.Time, dbRecord time.Time, weights models.Weights) float64 {

	if userClaim.IsZero() || dbRecord.IsZero() {
		return scoring.NoMatch
	}

	diff := math.Abs(userClaim.Sub(dbRecord).Hours()) / scoring.AvgAmountOfHoursInMonth
	toleranceMonths := scoring.ToleranceHours / scoring.AvgAmountOfHoursInMonth

	if diff <= scoring.IdealMatchofMonths+toleranceMonths {
		return weights.RegDate * scoring.IdealMatch
	}

	if diff <= scoring.PartialMatchofMonths+toleranceMonths {
		return weights.RegDate * scoring.PartialMatch
	}

	if diff-toleranceMonths >= scoring.OneYearofMonths {
		return -scoring.IdealMatch
	}

	return scoring.NoMatch
}
