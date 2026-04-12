package logic

import (
	"strings"
	"techsupport/core/internal/models"
	"techsupport/core/internal/scoring"
)

func CalculateScoreForRegCountry(userRegCountry, dbRegCountry string, weights models.Weights) float64 {

	Ucountry := strings.TrimSpace(userRegCountry)
	DCountry := strings.TrimSpace(dbRegCountry)

	if Ucountry == "" || DCountry == "" {
		return scoring.NoMatch
	}

	if Ucountry == "" && DCountry == "" {
		return scoring.NoMatch
	}

	if strings.EqualFold(Ucountry, DCountry) {
		return weights.RegCountry * scoring.IdealMatch
	}
	return scoring.NoMatch
}

func CalculateScoreForRegCity(userRegCity, dbRegCity string, weights models.Weights) float64 {

	uCity := strings.TrimSpace(userRegCity)
	dCity := strings.TrimSpace(dbRegCity)

	if uCity == "" || dCity == "" {
		return scoring.NoMatch
	}

	if uCity == "" && dCity == "" {
		return scoring.NoMatch
	}

	if strings.EqualFold(uCity, dCity) {
		return weights.RegCity * scoring.IdealMatch
	}

	return scoring.NoMatch
}
