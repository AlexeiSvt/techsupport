package transactions

import (
	"techsupport/internal/models"
	"techsupport/internal/scoring"
	"time"
)

func isRegionOrDeviceKnown(tx models.Transaction, history models.UserHistory) bool {
	// Оптимизировано: проверяем окна по очереди без append (экономим память)
	for _, session := range history.FirstWindow {
		if session.Country == tx.Country || session.DeviceID == tx.DeviceID {
			return true
		}
	}
	for _, session := range history.LastWindow {
		if session.Country == tx.Country || session.DeviceID == tx.DeviceID {
			return true
		}
	}
	return false
}

func calculateWindowScore(tx models.Transaction, history []models.Session) float64 {
	if len(history) == 0 {
		return scoring.NoMatch
	}

	var maxScore float64
	for _, session := range history {
		score := 0.0
		if session.City == tx.City { score += scoring.CityScore }
		if session.Country == tx.Country { score += scoring.CountryScore }
		if session.DeviceID == tx.DeviceID { score += scoring.DeviceScore }
		if session.SessionIP == tx.IP { score += scoring.IPScore }
		
		if score > maxScore {
			maxScore = score
		}
	}

	if maxScore < scoring.MinScoreForPartialMatch {
		return scoring.NoMatch
	}

	const maxPossibleScore = scoring.CityScore + scoring.CountryScore + scoring.DeviceScore + scoring.IPScore
	return maxScore / maxPossibleScore
}

func isSuddenHighDonation(tx models.Transaction, history models.UserHistory) bool {
	if len(history.AllPayments) == 0 {
		return tx.Amount >= scoring.FirstDonationThreshold
	}

	var total float64
	for _, p := range history.AllPayments {
		total += p.Amount
	}
	avg := total / float64(len(history.AllPayments))
	return tx.Amount > avg*scoring.SuddenMultiplier
}

func isHighFrequencyTransaction(allPayments []models.Transaction, current models.Transaction) bool {
	currentTime, err := time.Parse(time.RFC3339, current.Timestamp)
	if err != nil {
		return false
	}
	
	minInterval := time.Duration(scoring.MinIntervalHours * float64(time.Hour))
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