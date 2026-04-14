package engine

import (
	"fmt"
	"techsupport/core/internal/constants"
	"techsupport/core/internal/models"
	logPkg "techsupport/log/pkg"
)

type DeviceBruteforcePenaltyCalculator struct {
	Log logPkg.Logger
}

func (c *DeviceBruteforcePenaltyCalculator) Calculate(user models.UserData, db models.DBRecord, weights models.Weights) models.CalcResult {
	penalty := 0.0
	comment := "No penalties"

	if c.Log != nil {
		c.Log.Debugw("starting security penalty calculation",
			"user_devices_count", len(user.UserClaim.Devices),
			"db_devices_count", len(db.Devices),
		)
	}

	if len(user.UserClaim.Devices) > len(db.Devices)+constants.NewDevicesThreshold {
		penalty += constants.BruteforcePenalty
		comment = "New devices threshold exceeded"
		
		if c.Log != nil {
			c.Log.Warnw("device bruteforce penalty applied",
				"u_count", len(user.UserClaim.Devices),
				"db_count", len(db.Devices),
				"penalty", constants.BruteforcePenalty,
			)
		}
	}

	if user.UserClaim.IPInfo != nil {
		ipPenalty := user.UserClaim.IPInfo.GetPenaltyScore(c.Log)
		
		if ipPenalty > 0 {
			penalty += ipPenalty
			comment += fmt.Sprintf(" | IP Penalty: %.1f", ipPenalty)
		}
	} else {
		if c.Log != nil {
			c.Log.Warnw("IPInfo is missing during penalty calculation - potential data collection skip")
		}
	}

	if penalty > constants.FullPenalty {
		if c.Log != nil {
			c.Log.Debugw("capping penalty to max value", "raw_penalty", penalty, "max", constants.FullPenalty)
		}
		penalty = constants.FullPenalty
	}

	status := "clear"
	if penalty > 0 {
		status = "penalty"
		if c.Log != nil {
			c.Log.Infow("security penalty final result", "total_penalty", penalty, "status", status)
		}
	}

	return models.CalcResult{
		Name:    "Security Penalty",
		Code:    "penalty_score",
		Value:   penalty,
		Weight:  1.0, 
		Result:  penalty,
		Status:  status,
		Comment: comment,
	}
}