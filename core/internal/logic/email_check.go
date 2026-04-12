package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
)

var _ pkg.ScoreCalculator = (*FirstEmailCalculator)(nil)

type FirstEmailCalculator struct{}

func (c FirstEmailCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkEmail(user.UserClaim.FirstEmail, db.FirstEmail)

	return models.CalcResult{
		Name:    "First Email Match",
		Code:    "first_email",
		Value:   res.Value,
		Weight:  weights.FirstEmail,
		Result:  res.Value * weights.FirstEmail,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func checkEmail(userEmail, dbEmail string) rawCheckResult {
	if userEmail == "" || dbEmail == "" {
		return rawCheckResult{Value: constants.NoMatch, Status: "no_data", Comment: "One or both emails are empty"}
	}

	if strings.EqualFold(userEmail, dbEmail) {
		return rawCheckResult{Value: constants.IdealMatch, Status: "match", Comment: "Full case-insensitive match"}
	}

	return rawCheckResult{Value: constants.NoMatch, Status: "no_match", Comment: "Emails do not match"}
}
