package engine

import (
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
)

type DeviceBruteforcePenaltyCalculator struct{}

func (c DeviceBruteforcePenaltyCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
	score := 0.0

	if len(user.Devices) > len(db.Devices)+constants.NewDevicesThreshold {
		score += constants.BruteforcePenalty
	}

	if user.IPInfo != nil {
		score += user.IPInfo.GetPenaltyScore()
	}

	if score > constants.FullPenalty {
		score = constants.FullPenalty
	}

	return score
}