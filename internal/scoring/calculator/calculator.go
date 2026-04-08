package calculator

import (
    "techsupport/internal/models"
    "techsupport/internal/scoring/logic"
    "techsupport/internal/scoring/logic/transactions"
)

// ScoreCalculator теперь принимает UserData, чтобы соответствовать функции CalculateScore
type ScoreCalculator interface {
    Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64
}

type RegDateCalculator struct{}
func (c RegDateCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateDeltaOfCreationAge(user.UserClaim.RegDate, db.RegDate, weights)
}

type RegCountryCalculator struct{}
func (c RegCountryCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForRegCountry(user.UserClaim.RegCountry, db.RegCountry, weights)
}

type RegCityCalculator struct{}
func (c RegCityCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForRegCity(user.UserClaim.RegCity, db.RegCity, weights)
}

type FirstEmailCalculator struct{}
func (c FirstEmailCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForFirstEmail(user.UserClaim.FirstEmail, db.FirstEmail, weights)
}

type FirstPhoneCalculator struct{}
func (c FirstPhoneCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForFirstPhone(user.UserClaim.Phone, db.Phone, weights)
}

type FirstDeviceCalculator struct{}
func (c FirstDeviceCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForFirstDevice(user.UserClaim.FirstDevice, db.FirstDevice, weights)
}

type DevicesCalculator struct{}
func (c DevicesCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return logic.CalculateScoreForDevices(user.UserClaim.Devices, db, weights)
}

type FirstTransactionScoreCalculator struct{}
func (c FirstTransactionScoreCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) float64 {
    return transactions.CalculateFirstTransactionScore(db, user.UserClaim, weights)
}

func CalculateScore(user models.UserData, db models.DBRecord, weights models.Weights, calculators []ScoreCalculator) float64 {
    total := 0.0
    for _, calc := range calculators {
        total += calc.Calculate(user, db, weights)
    }
    return total
}