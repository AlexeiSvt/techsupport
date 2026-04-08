package transactions

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring"
)

func CalculateFirstTransactionScore(dbRecord models.DBRecord, userClaim models.UserClaim, weights models.Weights) float64 {
	if weights.FirstTransaction <= 0 {
		return scoring.NoMatch
	}

	tx := userClaim.FirstTransaction
	var baseScore float64

	if score := calculateWindowScore(tx, dbRecord.UserHistory.FirstWindow); score > 0 {
		baseScore = score
	}
	
	if score := calculateWindowScore(tx, dbRecord.UserHistory.LastWindow); score > 0 {
		if baseScore == 0 || score > baseScore {
			baseScore = score
		}
	}

	if baseScore == 0 {
		if isRegionOrDeviceKnown(tx, dbRecord.UserHistory) {
			baseScore = scoring.PartialMatch
		} else {
			return scoring.NoMatch
		}
	}

	if isHighFrequencyTransaction(dbRecord.UserHistory.AllPayments, tx) {
		baseScore *= scoring.MostlyMatch
	}

	if !isRegionOrDeviceKnown(tx, dbRecord.UserHistory) {
		baseScore *= scoring.MostlyMatch
	}

	if isSuddenHighDonation(tx, dbRecord.UserHistory) {
		baseScore *= scoring.MostlyMatch
	}

	return capTransactionScore(baseScore*weights.FirstTransaction, weights.FirstTransaction)
}

func capTransactionScore(score float64, maxWeight float64) float64 {
	if score > maxWeight {
		return maxWeight
	}
	return score
}