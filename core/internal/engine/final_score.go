package engine

import (
	"fmt"
	"math"
	"techsupport/core/internal/ipchecker"
	"techsupport/core/internal/logic"
	"techsupport/core/internal/models"
	logPkg "techsupport/log/pkg"
	"time"
)

func CalculateFinalScore(log logPkg.Logger, input models.InputData) models.OutputData {
	if log != nil {
		log.Infow("processing final score calculation", "ip", input.UserData.IP)
	}

	ipPenaltyScore := 0.0
	ipInfo, err := ipchecker.GetIpInfo(log, input.UserData.IP)
	
	if err != nil {
		if log != nil {
			log.Warnw("failed to get IP info, proceeding without IP penalty", "err", err)
		}
	} else {
		input.UserData.UserClaim.IPInfo = ipInfo
		input.UserData.ASN = ipInfo.GetOperator()
		
		ipPenaltyScore = ipInfo.GetPenaltyScore(log)
	}

	weights := logic.GetWeights(input.DBRecord.IsDonator)
	calculators := NewScoringEngine(log)

	details, baseScore := CalculateAll(log, input.UserData, input.DBRecord, weights, calculators)

	penaltyCalc := DeviceBruteforcePenaltyCalculator{Log: log}
	pRes := penaltyCalc.Calculate(input.UserData, input.DBRecord, weights)
	details = append(details, pRes)

	totalPenalty := pRes.Result + ipPenaltyScore
	effectivePenalty := math.Max(0, math.Min(100, totalPenalty))
	
	survivalRate := (100.0 - effectivePenalty) / 100.0

	finalScore := 0.0
	if baseScore > 0 {
		finalScore = baseScore * survivalRate
	}

	finalScore = math.Floor(finalScore*100) / 100

	if log != nil {
		log.Infow("final calculation complete", 
			"base_score", baseScore, 
			"total_penalty", totalPenalty, 
			"survival_rate", survivalRate, 
			"final_score", finalScore,
		)
	}

	return models.OutputData{
		TicketID:       0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		FinaPercentage: fmt.Sprintf("%.2f%%", finalScore),
		Metrics: models.Metrics{
			Knowledge:    baseScore,
			PenaltyScore: effectivePenalty,
		},
		Details: details,
	}
}