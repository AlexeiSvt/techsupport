package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
)

var _ pkg.ScoreCalculator = (*RegCountryCalculator)(nil)
var _ pkg.ScoreCalculator = (*RegCityCalculator)(nil)

type RegCountryCalculator struct {
	Log logPkg.Logger
}

func (c *RegCountryCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	if c.Log != nil {
		c.Log.Debugw("calculating country match", "u_country", user.UserClaim.RegCountry, "db_country", db.RegCountry)
	}

	res := c.checkLocation(user.UserClaim.RegCountry, db.RegCountry, "Country")

	return models.CalcResult{
		Name:    "Registration Country Match",
		Code:    "reg_country",
		Value:   res.Value,
		Weight:  weights.RegCountry,
		Result:  res.Value * weights.RegCountry,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

type RegCityCalculator struct {
	Log logPkg.Logger
}

func (c *RegCityCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	if c.Log != nil {
		c.Log.Debugw("calculating city match", "u_city", user.UserClaim.RegCity, "db_city", db.RegCity)
	}

	res := c.checkLocation(user.UserClaim.RegCity, db.RegCity, "City")

	return models.CalcResult{
		Name:    "Registration City Match",
		Code:    "reg_city",
		Value:   res.Value,
		Weight:  weights.RegCity,
		Result:  res.Value * weights.RegCity,
		Status:  res.Status,
		Comment: res.Comment,
	}
}

func (c *RegCountryCalculator) checkLocation(userVal, dbVal, label string) rawCheckResult {
	return checkLocationGeneric(c.Log, userVal, dbVal, label)
}

func (c *RegCityCalculator) checkLocation(userVal, dbVal, label string) rawCheckResult {
	return checkLocationGeneric(c.Log, userVal, dbVal, label)
}

func checkLocationGeneric(log logPkg.Logger, userVal, dbVal, label string) rawCheckResult {
	u := strings.TrimSpace(userVal)
	d := strings.TrimSpace(dbVal)

	if u == "" || d == "" {
		if log != nil {
			log.Warnw(errors.ErrEmptyLocationData.Error(), "label", label, "is_user_empty", u == "", "is_db_empty", d == "")
		}
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  errors.StatusNoData,
			Comment: label + " data is missing",
		}
	}

	if strings.EqualFold(u, d) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  errors.StatusMatch,
			Comment: label + " matches exactly",
		}
	}

	if log != nil {
		log.Debugw(errors.ErrLocationMismatch.Error(), "label", label, "u_val", u, "d_val", d)
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  errors.StatusNoMatch,
		Comment: label + " mismatch",
	}
}