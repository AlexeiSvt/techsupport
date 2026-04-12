package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
)

var _ pkg.ScoreCalculator = (*FirstPhoneCalculator)(nil)

type FirstPhoneCalculator struct{}

func (c FirstPhoneCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkPhone(user.UserClaim.Phone, db.Phone)
	
	return models.CalcResult{
		Name:    "First Phone Match",
		Code:    "first_phone",
		Value:   res.Value,
		Weight:  weights.Phone,
		Result:  res.Value * weights.Phone,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func checkPhone(userPhone, dbPhone string) rawCheckResult {
	if userPhone == "" || dbPhone == "" {
		return rawCheckResult{Value: constants.NoMatch, Status: "no_data", Comment: "Missing phone data"}
	}

	uPhone := CleanPhoneNumber(userPhone)
	dPhone := CleanPhoneNumber(dbPhone)

	if uPhone == dPhone {
		return rawCheckResult{Value: constants.IdealMatch, Status: "match", Comment: "Full phone number match"}
	}

	if len(uPhone) < constants.MinLen || len(dPhone) < constants.MinLen {
		return rawCheckResult{Value: constants.NoMatch, Status: "no_match", Comment: "Phone length below minimum"}
	}

	checkLen := min(min(len(dPhone), len(uPhone)), 10)

	if uPhone[len(uPhone)-checkLen:] == dPhone[len(dPhone)-checkLen:] {
		return rawCheckResult{Value: constants.PartialMatch, Status: "partial", Comment: "Matched by last 10 digits"}
	}

	return rawCheckResult{Value: constants.NoMatch, Status: "no_match", Comment: "Numbers do not match"}
}

func CleanPhoneNumber(phone string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
}