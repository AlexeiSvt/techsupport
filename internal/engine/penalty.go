package engine

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring"
)

type DeviceBruteforcePenaltyCalculator struct{}

func (c DeviceBruteforcePenaltyCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
	score := 0.0

	if len(user.Devices) > len(db.Devices)+scoring.NewDevicesThreshold {
		score += scoring.BruteforcePenalty
	}

	if user.IPInfo != nil {
		score += user.IPInfo.GetPenaltyScore()
	}

	if score > scoring.FullPenalty {
		score = scoring.FullPenalty
	}

	return score
}