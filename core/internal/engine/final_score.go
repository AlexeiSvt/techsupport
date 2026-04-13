package engine

import (
    "fmt"
    "math"
    "techsupport/core/internal/ipchecker"
    "techsupport/core/internal/logic"
    "techsupport/core/internal/models"
    "time"
)

func CalculateFinalScore(input models.InputData) models.OutputData {
    ipInfo, err := ipchecker.GetIpInfo(input.UserData.IP)
    
    ipPenaltyScore := 0.0
    
    if err == nil {
        input.UserData.UserClaim.IPInfo = ipInfo
        
  
        input.UserData.ASN = ipInfo.GetOperator()
        
        ipPenaltyScore = ipInfo.GetPenaltyScore()
    }

    weights := logic.GetWeights(input.DBRecord.IsDonator)
    calculators := NewScoringEngine()

    details, baseScore := CalculateAll(input.UserData, input.DBRecord, weights, calculators)

    penaltyCalc := DeviceBruteforcePenaltyCalculator{}
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