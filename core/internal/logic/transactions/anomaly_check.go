package transactions

import (
	"techsupport/core/internal/models"
	"techsupport/core/internal/constants"
	"time"
)

func isRegionAndDeviceKnown(tx models.Transaction, history models.UserHistory) bool {
    for _, session := range history.FirstWindow {
        if session.DeviceID == tx.DeviceID && (session.Country == tx.Country || session.City == tx.City) {
            return true
        }
    }

    for _, session := range history.LastWindow {
        if session.DeviceID == tx.DeviceID && (session.Country == tx.Country || session.City == tx.City) {
            return true
        }
    }
    return false
}

func calculateWindowScore(tx models.Transaction, history []models.Session) float64 {
	if len(history) == 0 {
		return constants.NoMatch
	}

	var maxScore float64
	for _, session := range history {
		score := 0.0
		if session.City == tx.City { score += constants.CityScore }
		if session.Country == tx.Country { score += constants.CountryScore }
		if session.DeviceID == tx.DeviceID { score += constants.DeviceScore }
		if session.SessionIP == tx.IP { score += constants.IPScore }
		
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore < constants.MinScoreForPartialMatch {
		return constants.NoMatch
	}

	const maxPossibleScore = constants.CityScore + constants.CountryScore + constants.DeviceScore + constants.IPScore
	return maxScore / maxPossibleScore
}

func isSuddenHighDonation(tx models.Transaction, history models.UserHistory) bool {
	if len(history.AllPayments) == 0 {
		return tx.Amount >= constants.FirstDonationThreshold
	}

	var total float64
	for _, p := range history.AllPayments {
		total += p.Amount
	}
	avg := total / float64(len(history.AllPayments))
	return tx.Amount > avg*constants.SuddenMultiplier
}

func isHighFrequencyTransaction(allPayments []models.Transaction, current models.Transaction) bool {
	currentTime, err := time.Parse(time.RFC3339, current.Timestamp)
	if err != nil {
		return false
	}
	
	minInterval := time.Duration(constants.MinIntervalHours * float64(time.Hour))
	for _, tx := range allPayments {
		txTime, err := time.Parse(time.RFC3339, tx.Timestamp)
		if err != nil {
			continue
		}
		
		diff := currentTime.Sub(txTime)
		if diff < 0 { diff = -diff }
		
		if diff < minInterval {
			return true
		}
	}
	return false
}