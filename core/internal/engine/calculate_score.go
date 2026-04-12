package engine

import (
	"techsupport/core/internal/models"
	"techsupport/core/pkg"
)

func CalculateScore(user models.UserData, db models.DBRecord, weights models.Weights, calculators []pkg.ScoreCalculator) float64 { 
    total := 0.0
    
    for _, calc := range calculators {
        res := calc.Calculate(user, db, weights)
        
        total += res.Result
    }
    
    return total
}