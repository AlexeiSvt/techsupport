package engine

import (
	"math"
	"techsupport/core/internal/ipchecker"
	"techsupport/core/internal/models"
	"techsupport/core/internal/logic"
)

func CalculateFinalScore(input models.InputData) float64 {

    ipInfo, err := ipchecker.GetIpInfo(input.UserData.IP)
    
    var ipPenalty float64 = 100.0
    if err == nil {
        ipPenalty = ipInfo.GetPenaltyScore()
        input.UserData.ASN = ipInfo.ASN.Org 
    }

    if ipPenalty >= 100 {
        return 0.0
    }

    penaltyCalc := DeviceBruteforcePenaltyCalculator{}
    weights := logic.GetWeights(input.DBRecord.IsDonator)
    devicePenalty := penaltyCalc.Calculate(input.UserData.UserClaim, input.DBRecord, weights)

    devicePenalty = math.Max(0, math.Min(100, devicePenalty))

    baseScore := CalculateScoreForClaim(input)
    if baseScore <= 0 {
        return 0.0
    }

    survivalRate := (100.0 - devicePenalty) / 100.0

    finalScore := baseScore * survivalRate

    return float64(math.Floor(finalScore*100) / 100)
}