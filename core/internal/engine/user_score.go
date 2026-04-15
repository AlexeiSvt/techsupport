package engine

import (
    "techsupport/core/internal/logic"
    "techsupport/core/pkg/models"
    "techsupport/core/internal/logic/transactions" 
    "techsupport/core/pkg"
    logPkg "techsupport/log/pkg"
)

func NewScoringEngine(log logPkg.Logger) []pkg.ScoreCalculator {
    txValidator := &transactions.TxValidator{Log: log}

    return []pkg.ScoreCalculator{
        &logic.RegDateCalculator{Log: log},
        &logic.RegCountryCalculator{Log: log},
        &logic.RegCityCalculator{Log: log},
        &logic.FirstEmailCalculator{Log: log},
        &logic.FirstPhoneCalculator{Log: log},
        &logic.FirstDeviceCalculator{Log: log},
        &logic.DevicesCalculator{Log: log},
        &transactions.FirstTransactionScoreCalculator{
            Log:       log,
            Validator: txValidator,
        },
        &DeviceBruteforcePenaltyCalculator{Log: log},
    }
}

func CalculateAll(log logPkg.Logger, user models.UserData, db models.DBRecord, weights models.Weights, calculators []pkg.ScoreCalculator) ([]models.CalcResult, float64) {
    var details []models.CalcResult
    var total float64

    if log != nil {
        log.Infow("starting full scoring session", 
            "is_donator", weights.FirstTransaction > 0,
        )
    }

    for _, calc := range calculators {
        res := calc.Calculate(user, db, weights)
        
        total += res.Result 
        details = append(details, res)
    }

    if log != nil {
        log.Infow("scoring session completed", 
            "total_score", total, 
            "results_count", len(details),
        )
    }

    return details, total
}