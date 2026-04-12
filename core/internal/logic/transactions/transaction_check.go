package transactions

import (
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
)

func CalculateFirstTransactionScore(dbRecord models.DBRecord, userClaim models.UserClaim, weights models.Weights) float64 {
    if weights.FirstTransaction <= 0 {
        return constants.NoMatch
    }

    tx := userClaim.FirstTransaction
    var baseScore float64
    anomalyCount := 0

    if score := calculateWindowScore(tx, dbRecord.UserHistory.FirstWindow); score > 0 {
        baseScore = score
    }
    
    if score := calculateWindowScore(tx, dbRecord.UserHistory.LastWindow); score > 0 {
        if baseScore == 0 || score > baseScore {
            baseScore = score
        }
    }

    if baseScore == 0 {
        if isRegionAndDeviceKnown(tx, dbRecord.UserHistory) {
            baseScore = constants.PartialMatch
        } else {
            return constants.NoMatch 
        }
    }

    if isHighFrequencyTransaction(dbRecord.UserHistory.AllPayments, tx) {
        anomalyCount++
        baseScore *= constants.MostlyMatch
    }

    if isSuddenHighDonation(tx, dbRecord.UserHistory) {
        anomalyCount++
        baseScore *= constants.MostlyMatch
    }

    if !isRegionAndDeviceKnown(tx, dbRecord.UserHistory) {
        anomalyCount++
        baseScore *= constants.MostlyMatch
    }

    if anomalyCount >= 2 {
        return constants.NoMatch 
    }

    return capTransactionScore(baseScore * weights.FirstTransaction, weights.FirstTransaction)
}

func capTransactionScore(score float64, maxWeight float64) float64 {
	if score > maxWeight {
		return maxWeight
	}
	return score
}