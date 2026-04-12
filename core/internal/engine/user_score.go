package engine

import (
    "techsupport/core/internal/logic"
    "techsupport/core/internal/logic/transactions"
    "techsupport/core/internal/models"
    "techsupport/core/pkg"
)

func NewScoringEngine() []pkg.ScoreCalculator {
    return []pkg.ScoreCalculator{
        logic.RegDateCalculator{},
        logic.RegCountryCalculator{},
        logic.RegCityCalculator{},
        logic.FirstEmailCalculator{},
        logic.FirstPhoneCalculator{},
        logic.FirstDeviceCalculator{},
        logic.DevicesCalculator{},
        transactions.FirstTransactionScoreCalculator{},
    }
}

func CalculateAll(user models.UserData, db models.DBRecord, weights models.Weights, calculators []pkg.ScoreCalculator) ([]models.CalcResult, float64) {
    var details []models.CalcResult
    var total float64

    for _, calc := range calculators {
        res := calc.Calculate(user, db, weights)
        
        total += res.Result 
        details = append(details, res)
    }

    return details, total
}