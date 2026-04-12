package logic

import (
	"strings"
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
)

func CalculateScoreForRegCountry(userRegCountry, dbRegCountry string, weights models.Weights) float64 {

	Ucountry := strings.TrimSpace(userRegCountry)
	DCountry := strings.TrimSpace(dbRegCountry)

	if Ucountry == "" || DCountry == "" {
		return constants.NoMatch
	}

	if Ucountry == "" && DCountry == "" {
		return constants.NoMatch
	}

	if strings.EqualFold(Ucountry, DCountry) {
		return weights.RegCountry * constants.IdealMatch
	}
	return constants.NoMatch
}

func CalculateScoreForRegCity(userRegCity, dbRegCity string, weights models.Weights) float64 {

	uCity := strings.TrimSpace(userRegCity)
	dCity := strings.TrimSpace(dbRegCity)

	if uCity == "" || dCity == "" {
		return constants.NoMatch
	}

	if uCity == "" && dCity == "" {
		return constants.NoMatch
	}

	if strings.EqualFold(uCity, dCity) {
		return weights.RegCity * constants.IdealMatch
	}

	return constants.NoMatch
}
