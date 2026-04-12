package engine

import (
	"techsupport/core/internal/models"
	"techsupport/core/internal/logic"

	"techsupport/core/pkg"
)

func CalculateScoreForClaim(input models.InputData) float64 {
	weights := logic.GetWeights(input.DBRecord.IsDonator)

	calculators := []pkg.ScoreCalculator{
		pkg.RegDateCalculator{},
		pkg.RegCountryCalculator{},
		pkg.RegCityCalculator{},
		pkg.FirstEmailCalculator{},
		pkg.FirstPhoneCalculator{},
		pkg.FirstDeviceCalculator{},
		pkg.DevicesCalculator{},
		pkg.FirstTransactionScoreCalculator{},
	}

	totalScore := pkg.CalculateScore(
		input.UserData,
		input.DBRecord,
		weights,
		calculators,
	)

	return totalScore
}