package calculator

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring/logic"
)

type ScoreCalculator interface {
    Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64
}

type RegDateCalculator struct{}
func (c RegDateCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateDeltaOfCreationAge(user.RegDate, db.RegDate, weights)
}

type RegCountryCalculator struct{}
func (c RegCountryCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForRegCountry(user.RegCountry, db.RegCountry, weights)
}

type RegCityCalculator struct{}
func (c RegCityCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForRegCity(user.RegCity, db.RegCity, weights)
}

type FirstEmailCalculator struct{}
func (c FirstEmailCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForFirstEmail(user.FirstEmail, db.FirstEmail, weights)
}

type FirstPhoneCalculator struct{}
func (c FirstPhoneCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForFirstPhone(user.Phone, db.Phone, weights)
}

type FirstDeviceCalculator struct{}
func (c FirstDeviceCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForFirstDevice(user.FirstDevice, db.FirstDevice, weights)
}

type DevicesCalculator struct{}
func (c DevicesCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForDevices(user.Devices, db, weights)
}

type FirstTransactionScoreCalculator struct{}
func (c FirstTransactionScoreCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    return logic.FirstTransactionScoreCalculator(db, user, weights)
}

func CalculateScore(user models.UserClaim, db models.DBRecord, weights models.Weights, calculators []ScoreCalculator) float64 {
    total := 0.0
    for _, calc := range calculators {
        total += calc.Calculate(user, db, weights)
    }
    return total
}