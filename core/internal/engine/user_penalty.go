package engine

import (
    "techsupport/core/internal/models"
    "techsupport/core/internal/constants"
    "fmt"
)

type DeviceBruteforcePenaltyCalculator struct{}

func (c DeviceBruteforcePenaltyCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
    penalty := 0.0
    comment := "No penalties"

    if len(user.UserClaim.Devices) > len(db.Devices)+constants.NewDevicesThreshold {
        penalty += constants.BruteforcePenalty
        comment = "New devices threshold exceeded"
    }

    if user.UserClaim.IPInfo != nil {
        ipPenalty := user.UserClaim.IPInfo.GetPenaltyScore()
        penalty += ipPenalty
        comment += fmt.Sprintf(" | IP Penalty: %.1f", ipPenalty)
    }

    if penalty > constants.FullPenalty {
        penalty = constants.FullPenalty
    }

    return models.CalcResult{
        Name:    "Security Penalty",
        Code:    "penalty_score",
        Value:   penalty,
        Weight:  1.0, 
        Result:  penalty,
        Status:  "penalty",
        Comment: comment,
    }
}