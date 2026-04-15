package engine

import (
	"techsupport/core/pkg/models"
	"techsupport/core/pkg"
	logPkg "techsupport/log/pkg"
)

func CalculateScore(log logPkg.Logger, user models.UserData, db models.DBRecord, weights models.Weights, calculators []pkg.ScoreCalculator) float64 {
	total := 0.0

	if log != nil {
		log.Debugw("starting fast score calculation",
			"calculators_count", len(calculators),
		)
	}

	for _, calc := range calculators {
		res := calc.Calculate(user, db, weights)
		total += res.Result
	}

	// Проверка перед выходом
	if log != nil {
		log.Infow("fast score calculation finished",
			"total_score", total,
		)
	}

	return total
}