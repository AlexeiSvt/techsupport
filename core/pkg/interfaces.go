package pkg

import "techsupport/core/pkg/models"

type ScoreCalculator interface {
	Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult
}

type Engine interface {
	Run(input models.InputData) models.OutputData
}