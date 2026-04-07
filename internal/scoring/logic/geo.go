package logic

import (
	"strings"
	"techsupport/internal/models"
	"techsupport/internal/scoring"
)

func CalculateScoreForRegCountry(userRegCountry, dbRegCountry string, weights models.Weights) float64 {
	
	if userRegCountry == "" && dbRegCountry == "" {
		return scoring.NoMatch
	}

	if userRegCountry == "" || dbRegCountry == "" {
		return scoring.NoMatch
	}

	if strings.EqualFold(userRegCountry, dbRegCountry) {
		return weights.RegCountry * scoring.IdealMatch
	}
	return scoring.NoMatch
}

func CalculateScoreForRegCity(userRegCity, dbRegCity string, weights models.Weights) float64 {
	if userRegCity == "" && dbRegCity == "" {
		return scoring.NoMatch
	}

	if userRegCity == "" || dbRegCity == "" {
		return scoring.NoMatch
	}

	if strings.EqualFold(userRegCity, dbRegCity) {
		return weights.RegCity * scoring.IdealMatch
	}
	return scoring.NoMatch
}