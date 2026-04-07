package engine

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring"
)

type DeviceBruteforcePenaltyCalculator struct{}

func (c DeviceBruteforcePenaltyCalculator) Calculate(user models.UserClaim, db models.DBRecord, weights models.Weights) float64 {
    if len(user.Devices) > len(db.Devices) + scoring.NewDevicesThreshold {
        return scoring.BruteforcePenalty
    }
    return 0
}