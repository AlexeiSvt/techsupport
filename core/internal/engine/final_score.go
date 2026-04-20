// Package engine orchestrates the high-level scoring workflow.
package engine

import (
	"context"
	"fmt"
	"math"
	"time"

	"techsupport/core/internal/ipchecker"
	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// CalculateFinalScore is the main pipeline of the scoring system.
// It fetches external IP data, runs all logic calculators, applies security penalties,
// and produces the final weighted percentage.
func CalculateFinalScore(ctx context.Context, log logPkg.Logger, input models.InputData) models.OutputData {
	if log != nil {
		log.Infow("processing final score calculation", "ip", input.UserData.IPInfo.IP)
	}

	// 1. IP Intelligence Phase
	ipPenaltyScore := 0.0
	// We use the context-aware version to ensure we don't hang on network calls.
	ipInfo, err := ipchecker.GetIpInfoWithContext(ctx, log, input.UserData.IPInfo.IP)
	
	if err != nil {
		if log != nil {
			log.Warnw("failed to get IP info, proceeding without IP penalty", "err", err)
		}
	} else {
		// Enrich user data with live IP intelligence.
		input.UserData.IPInfo = ipInfo
		input.UserData.IPInfo.ASN = ipInfo.ASN
		
		// Initial penalty based purely on IP reputation.
		ipPenaltyScore = ipInfo.GetPenaltyScore(log)
	}

	// 2. Core Scoring Phase
	// Weights are determined by the user's status (e.g., Donator vs Regular).
	weights := logic.GetWeights(input.DBRecord.IsDonator)
	calculators := NewScoringEngine(log)

	// Run all registered calculators to get the Base Score and full audit details.
	details, baseScore := CalculateAll(ctx, log, input.UserData, input.DBRecord, weights, calculators)

	// 3. Security Penalty Phase
	penaltyCalc := NewDeviceBruteforcePenaltyCalculator(log)
	pRes := penaltyCalc.Calculate(ctx, input.UserData, input.DBRecord, weights)
	details = append(details, pRes)

	// 4. Mathematical Synthesis
	// Total penalty is the sum of IP risks and behavioral (bruteforce) risks.
	totalPenalty := pRes.Result + ipPenaltyScore
	// Clamp the penalty between 0 and 100%.
	effectivePenalty := math.Max(0, math.Min(100, totalPenalty))
	
	// Survival Rate determines how much of the Base Score remains after penalties.
	// Example: 20% penalty means survival rate is 0.8.
	survivalRate := (100.0 - effectivePenalty) / 100.0

	finalScore := 0.0
	if baseScore > 0 {
		finalScore = baseScore * survivalRate
	}

	// Round to two decimal places for financial-grade precision.
	finalScore = math.Floor(finalScore*100) / 100

	if log != nil {
		log.Infow("final calculation complete", 
			"base_score", baseScore, 
			"total_penalty", totalPenalty, 
			"survival_rate", survivalRate, 
			"final_score", finalScore,
		)
	}

	// 5. Output Construction
	return models.OutputData{
		TicketID:       0, // Should be assigned by the caller or DB.
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