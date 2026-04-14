package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
)

var _ pkg.ScoreCalculator = (*FirstPhoneCalculator)(nil)

type FirstPhoneCalculator struct {
	Log logPkg.Logger
}

func (c *FirstPhoneCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	if c.Log != nil {
		c.Log.Debugw("calculating phone match score",
			"user_phone_raw", user.UserClaim.Phone,
			"db_phone_raw", db.Phone,
		)
	}

	res := c.checkPhone(user.UserClaim.Phone, db.Phone)

	calcRes := models.CalcResult{
		Name:    "First Phone Match",
		Code:    "first_phone",
		Value:   res.Value,
		Weight:  weights.Phone,
		Result:  res.Value * weights.Phone,
		Status:  res.Status,
		Comment: res.Comment,
	}

	if c.Log != nil {
		c.Log.Infow("phone score calculation finished",
			"status", calcRes.Status,
			"result", calcRes.Result,
		)
	}

	return calcRes
}

func (c *FirstPhoneCalculator) checkPhone(userPhone, dbPhone string) rawCheckResult {
	if userPhone == "" || dbPhone == "" {
		if c.Log != nil {
			c.Log.Warnw(errors.ErrEmptyPhoneData.Error(), "u_phone_empty", userPhone == "", "db_phone_empty", dbPhone == "")
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: "Missing phone data",
		}
	}

	uPhone := CleanPhoneNumber(userPhone)
	dPhone := CleanPhoneNumber(dbPhone)

	if uPhone == dPhone {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: "Full phone number match",
		}
	}

	if len(uPhone) < constants.MinLen || len(dPhone) < constants.MinLen {
		if c.Log != nil {
			c.Log.Debugw(errors.ErrPhoneTooShort.Error(), "u_len", len(uPhone), "d_len", len(dPhone))
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoMatch,
			Comment: "Phone length below minimum",
		}
	}

	// Calculate suffix match
	checkLen := 10
	if len(dPhone) < checkLen {
		checkLen = len(dPhone)
	}
	if len(uPhone) < checkLen {
		checkLen = len(uPhone)
	}

	if uPhone[len(uPhone)-checkLen:] == dPhone[len(dPhone)-checkLen:] {
		if c.Log != nil {
			c.Log.Debugw("partial phone match detected", "check_len", checkLen)
		}
		return rawCheckResult{
			Value:   constants.PartialMatch,
			Status:  errors.StatusPartial,
			Comment: "Matched by last digits",
		}
	}

	if c.Log != nil {
		c.Log.Debugw(errors.ErrPhoneNoMatch.Error(), "u_cleaned", uPhone, "d_cleaned", dPhone)
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: "Numbers do not match",
	}
}

func CleanPhoneNumber(phone string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)
}