package logic

import (
	"fmt"
	"math"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
	"time"
)

var _ pkg.ScoreCalculator = (*RegDateCalculator)(nil)

type RegDateCalculator struct {
	Log logPkg.Logger
}

func (c *RegDateCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	// Исправлено: добавлена проверка на nil
	if c.Log != nil {
		c.Log.Debugw("calculating registration date match",
			"user_reg_date", user.UserClaim.RegDate,
			"db_reg_date", db.RegDate,
		)
	}

	res := c.checkCreationAge(user.UserClaim.RegDate, db.RegDate)

	calcRes := models.CalcResult{
		Name:    "Registration Date Match",
		Code:    "reg_date",
		Value:   res.Value,
		Weight:  weights.RegDate,
		Result:  res.Value * weights.RegDate,
		Status:  res.Status,
		Comment: res.Comment,
	}

	// Исправлено: добавлена проверка на nil
	if c.Log != nil {
		c.Log.Infow("reg date score calculated",
			"status", calcRes.Status,
			"diff_value", calcRes.Value,
		)
	}

	return calcRes
}

func (c *RegDateCalculator) checkCreationAge(userClaim time.Time, dbRecord time.Time) rawCheckResult {
	if userClaim.IsZero() || dbRecord.IsZero() {
		// Исправлено: добавлена проверка на nil
		if c.Log != nil {
			c.Log.Warnw(errors.ErrEmptyRegDate.Error(),
				"user_date_zero", userClaim.IsZero(),
				"db_date_zero", dbRecord.IsZero(),
			)
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing registration date in claim or database",
		}
	}

	diffHours := math.Abs(userClaim.Sub(dbRecord).Hours())
	diffMonths := diffHours / constants.AvgAmountOfHoursInMonth
	toleranceMonths := constants.ToleranceHours / constants.AvgAmountOfHoursInMonth

	commentBase := fmt.Sprintf("Diff: %.1f months (tolerance: %.1f)", diffMonths, toleranceMonths)

	// Исправлено: добавлена проверка на nil
	if c.Log != nil {
		c.Log.Debugw("date difference calculated", "diff_months", diffMonths, "tolerance", toleranceMonths)
	}

	if diffMonths <= constants.IdealMatchofMonths+toleranceMonths {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: commentBase + " - Within ideal range",
		}
	}

	if diffMonths <= constants.PartialMatchofMonths+toleranceMonths {
		return rawCheckResult{
			Value:   constants.PartialMatch,
			Status:  errors.StatusPartial,
			Comment: commentBase + " - Within partial range",
		}
	}

	if diffMonths-toleranceMonths >= constants.OneYearofMonths {

		if c.Log != nil {
			c.Log.Warnw(errors.ErrRegDateAnomaly.Error(), "diff_months", diffMonths)
		}

		return rawCheckResult{
			Value:   -constants.IdealMatch,
			Status:  errors.StatusAnomaly,
			Comment: commentBase + " - High discrepancy (over 1 year)",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: commentBase + " - Outside acceptable ranges",
	}
}