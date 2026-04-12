package logic

import (
	"strings"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
)

var _ pkg.ScoreCalculator = (*RegCountryCalculator)(nil)
var _ pkg.ScoreCalculator = (*RegCityCalculator)(nil)

type RegCountryCalculator struct{}

func (c RegCountryCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkLocation(user.UserClaim.RegCountry, db.RegCountry, "Country")

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

// --- CITY CALCULATOR ---

type RegCityCalculator struct{}

func (c RegCityCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	res := checkLocation(user.UserClaim.RegCity, db.RegCity, "City")

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

// Универсальная функция для проверки текстовых локаций
func checkLocation(userVal, dbVal, label string) rawCheckResult {
	u := strings.TrimSpace(userVal)
	d := strings.TrimSpace(dbVal)

	if u == "" || d == "" {
		return rawCheckResult{
			Value:   constants.NoMatch,
			Status:  "no_data",
			Comment: label + " data is missing in claim or database",
		}
	}

	if strings.EqualFold(u, d) {
		return rawCheckResult{
			Value:   constants.IdealMatch,
			Status:  "match",
			Comment: label + " matches exactly (case-insensitive)",
		}
	}

	return rawCheckResult{
		Value:   constants.NoMatch,
		Status:  "no_match",
		Comment: label + " mismatch: " + u + " vs " + d,
	}
}