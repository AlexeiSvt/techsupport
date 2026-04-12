package logic

import (
	"math"
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
	"time"
)

func CalculateDeltaOfCreationAge(userClaim time.Time, dbRecord time.Time, weights models.Weights) float64 {

	if userClaim.IsZero() || dbRecord.IsZero() {
		return constants.NoMatch
	}

	diff := math.Abs(userClaim.Sub(dbRecord).Hours()) / constants.AvgAmountOfHoursInMonth
	toleranceMonths := constants.ToleranceHours / constants.AvgAmountOfHoursInMonth

	if diff <= constants.IdealMatchofMonths+toleranceMonths {
		return weights.RegDate * constants.IdealMatch
	}

	if diff <= constants.PartialMatchofMonths+toleranceMonths {
		return weights.RegDate * constants.PartialMatch
	}

	if diff-toleranceMonths >= constants.OneYearofMonths {
		return -constants.IdealMatch
	}

	return constants.NoMatch
}
