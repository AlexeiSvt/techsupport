package engine

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring/logic"

	interfaces "techsupport/internal/scoring/calculator"
)

func CalculateScoreForClaim(input models.InputData) float64 {
	weights := logic.GetWeights(input.DBRecord.IsDonator)

	calculators := []interfaces.ScoreCalculator{
		interfaces.RegDateCalculator{},
		interfaces.RegCountryCalculator{},
		interfaces.RegCityCalculator{},
		interfaces.FirstEmailCalculator{},
		interfaces.FirstPhoneCalculator{},
		interfaces.FirstDeviceCalculator{},
		interfaces.DevicesCalculator{},
		interfaces.FirstTransactionScoreCalculator{},
	}

	totalScore := interfaces.CalculateScore(
		input.UserData,
		input.DBRecord,
		weights,
		calculators,
	)

	return totalScore
}