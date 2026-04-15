package transactions

import (
	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
	"time"
)

type TxValidator struct {
	Log logPkg.Logger
}

func (v *TxValidator) IsRegionAndDeviceKnown(tx models.Transaction, history models.UserHistory) bool {
	check := func(sessions []models.Session, windowName string) bool {
		for _, s := range sessions {
			if s.DeviceID == tx.DeviceID && (s.Country == tx.Country || s.City == tx.City) {
				if v.Log != nil {
					v.Log.Debugw("match found in window", "window", windowName, "device", s.DeviceID)
				}
				return true
			}
		}
		return false
	}

	known := check(history.FirstWindow, "first") || check(history.LastWindow, "last")
	if !known && v.Log != nil {
		v.Log.Warnw("unknown region or device for transaction", "tx_id", tx.TransactionID, "device", tx.DeviceID)
	}
	return known
}

func (v *TxValidator) CalculateWindowScore(tx models.Transaction, history []models.Session) float64 {
	if len(history) == 0 {
		return 0
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
		if v.Log != nil {
			v.Log.Debugw("window score below threshold", "max_score", maxScore)
		}
		return 0
	}

	const maxPossibleScore = constants.CityScore + constants.CountryScore + constants.DeviceScore + constants.IPScore
	finalScore := maxScore / maxPossibleScore
	
	if v.Log != nil {
		v.Log.Debugw("window score calculated", "score", finalScore)
	}
	return finalScore
}

func (v *TxValidator) IsSuddenHighDonation(tx models.Transaction, history models.UserHistory) bool {
	if len(history.AllPayments) == 0 {
		isHigh := tx.Amount >= constants.FirstDonationThreshold
		if isHigh && v.Log != nil {
			v.Log.Warnw(errors.ErrSuddenHighAmount.Error(), "amount", tx.Amount, "type", "first_donation")
		}
		return isHigh
	}

	var total float64
	for _, p := range history.AllPayments {
		total += p.Amount
	}
	avg := total / float64(len(history.AllPayments))
	isSudden := tx.Amount > avg*constants.SuddenMultiplier

	if isSudden && v.Log != nil {
		v.Log.Warnw(errors.ErrSuddenHighAmount.Error(), "amount", tx.Amount, "avg", avg, "multiplier", constants.SuddenMultiplier)
	}
	return isSudden
}

func (v *TxValidator) IsHighFrequencyTransaction(allPayments []models.Transaction, current models.Transaction) bool {
	currentTime, err := time.Parse(time.RFC3339, current.Timestamp)
	if err != nil {
		if v.Log != nil {
			v.Log.Errorw(errors.ErrTxTimeParse.Error(), "raw_timestamp", current.Timestamp, "err", err)
		}
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
			if v.Log != nil {
				v.Log.Warnw(errors.ErrHighFreqTx.Error(), "diff_sec", diff.Seconds(), "limit_sec", minInterval.Seconds())
			}
			return true
		}
	}
	return false
}