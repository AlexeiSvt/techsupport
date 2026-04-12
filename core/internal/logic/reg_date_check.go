package logic

import (
	"fmt"
	"math"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
	"time"
)

var _ pkg.ScoreCalculator = (*RegDateCalculator)(nil)

type RegDateCalculator struct{}

func (c RegDateCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkCreationAge(user.UserClaim.RegDate, db.RegDate)

	return models.CalcResult{
		Name:    "Registration Date Match",
		Code:    "reg_date",
		Value:   res.Value,
		Weight:  weights.RegDate,
		Result:  res.Value * weights.RegDate,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func checkCreationAge(userClaim time.Time, dbRecord time.Time) rawCheckResult {
	if userClaim.IsZero() || dbRecord.IsZero() {
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  "no_data",
			Comment: "Missing registration date in claim or database",
		}
	}

	diffHours := math.Abs(userClaim.Sub(dbRecord).Hours())
	diffMonths := diffHours / constants.AvgAmountOfHoursInMonth
	toleranceMonths := constants.ToleranceHours / constants.AvgAmountOfHoursInMonth

	commentBase := fmt.Sprintf("Diff: %.1f months (tolerance: %.1f)", diffMonths, toleranceMonths)

	if diffMonths <= constants.IdealMatchofMonths+toleranceMonths {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  "match",
			Comment: commentBase + " - Within ideal range",
		}
	}

	if diffMonths <= constants.PartialMatchofMonths+toleranceMonths {
		return rawCheckResult{
			Value:   constants.PartialMatch,
			Status:  "partial",
			Comment: commentBase + " - Within partial range",
		}
	}

	if diffMonths-toleranceMonths >= constants.OneYearofMonths {
		return rawCheckResult{
			Value:   -constants.IdealMatch,
			Status:  "anomaly",
			Comment: commentBase + " - High discrepancy (over 1 year)",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  "no_match",
		Comment: commentBase + " - Outside acceptable ranges",
	}
}
