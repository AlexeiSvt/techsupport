package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
)

var _ pkg.ScoreCalculator = (*FirstEmailCalculator)(nil)

type FirstEmailCalculator struct {
	Log logPkg.Logger
}

func (c *FirstEmailCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	userEmail := user.UserClaim.FirstEmail
	dbEmail := db.FirstEmail

	if c.Log != nil {
		c.Log.Debugw("starting email match calculation", "user_email", userEmail, "db_email", dbEmail)
	}

	res := c.checkEmail(userEmail, dbEmail)

	calcRes := models.CalcResult{
		Name:    "First Email Match",
		Code:    "first_email",
		Value:   res.Value,
		Weight:  weights.FirstEmail,
		Result:  res.Value * weights.FirstEmail,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if c.Log != nil {
		c.Log.Infow("email score calculation finished",
			"status", calcRes.Status,
			"score_impact", calcRes.Result,
		)
	}

	return calcRes
}

func (c *FirstEmailCalculator) checkEmail(userEmail, dbEmail string) rawCheckResult {
	if userEmail == "" || dbEmail == "" {
		if c.Log != nil {
			c.Log.Warnw(errors.ErrEmptyEmailData.Error(), "user_email_empty", userEmail == "", "db_email_empty", dbEmail == "")
		}

		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing email data for comparison",
		}
	}

	if strings.EqualFold(userEmail, dbEmail) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: "Full case-insensitive match",
		}
	}

	if c.Log != nil {
		c.Log.Debugw(errors.ErrEmailNoMatch.Error(), "u_email", userEmail, "db_email", dbEmail)
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: "Emails do not match",
	}
}