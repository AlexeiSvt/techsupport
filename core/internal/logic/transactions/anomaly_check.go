// Package transactions provides validation logic for financial activities.
// It analyzes user history, geographic consistency, and transaction patterns.
package transactions

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"techsupport/core/internal/constants"
	"techsupport/core/internal/errors"
	"techsupport/core/pkg/models"
	logPkg "techsupport/log/pkg"
)

// TxValidator performs risk assessment on incoming transactions.
// It is thread-safe and can be used concurrently across multiple requests.
type TxValidator struct {
	mu  sync.RWMutex // Protects the logger instance during runtime updates.
	log logPkg.Logger

	totalChecks uint64 // Atomic counter for total validation attempts.
	fraudAlerts uint64 // Atomic counter for detected anomalies.
}

// NewTxValidator creates a new instance of TxValidator with the provided logger.
func NewTxValidator(logger logPkg.Logger) *TxValidator {
	return &TxValidator{
		log: logger,
	}
}

// IsRegionAndDeviceKnown checks if the transaction's device and location 
// have been previously seen in the user's history windows.
func (v *TxValidator) IsRegionAndDeviceKnown(ctx context.Context, tx models.Transaction, history models.UserHistory) bool {
	atomic.AddUint64(&v.totalChecks, 1)

	v.mu.RLock()
	logger := v.log
	v.mu.RUnlock()

	// Internal helper to scan specific session windows.
	check := func(sessions []models.Session, windowName string) bool {
		for _, s := range sessions {
			// A match is found if the device matches AND (City or Country matches).
			if s.DeviceID == tx.DeviceID && (s.Country == tx.Country || s.City == tx.City) {
				if logger != nil {
					logger.Debugw("match found in window", "window", windowName, "device", s.DeviceID)
				}
				return true
			}
		}
		return false
	}

	known := check(history.FirstWindow, "first") || check(history.LastWindow, "last")
	
	if !known {
		atomic.AddUint64(&v.fraudAlerts, 1)
		if logger != nil {
			logger.Warnw("unknown region or device for transaction", 
				"tx_id", tx.TransactionID, 
				"device", tx.DeviceID,
			)
		}
	}
	return known
}

// CalculateWindowScore assigns a numerical confidence score based on session history.
// It returns a normalized value [0.0, 1.0] representing the quality of the match.
func (v *TxValidator) CalculateWindowScore(ctx context.Context, tx models.Transaction, history []models.Session) float64 {
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
		return 0
	}

	const maxPossibleScore = constants.CityScore + constants.CountryScore + constants.DeviceScore + constants.IPScore
	return maxScore / maxPossibleScore
}

// IsSuddenHighDonation detects if the transaction amount is significantly 
// higher than the user's historical average or initial threshold.
func (v *TxValidator) IsSuddenHighDonation(ctx context.Context, tx models.Transaction, history models.UserHistory) bool {
	v.mu.RLock()
	logger := v.log
	v.mu.RUnlock()

	// Scenario 1: No payment history (First donation check).
	if len(history.AllPayments) == 0 {
		isHigh := tx.Amount >= constants.FirstDonationThreshold
		if isHigh && logger != nil {
			logger.Warnw(errors.ErrSuddenHighAmount.Error(), "amount", tx.Amount, "type", "first_donation")
		}
		return isHigh
	}

	// Scenario 2: Compare against historical average.
	var total float64
	for _, p := range history.AllPayments {
		total += p.Amount
	}
	avg := total / float64(len(history.AllPayments))
	isSudden := tx.Amount > avg*constants.SuddenMultiplier

	if isSudden && logger != nil {
		logger.Warnw(errors.ErrSuddenHighAmount.Error(), 
			"amount", tx.Amount, 
			"avg", avg, 
			"multiplier", constants.SuddenMultiplier,
		)
	}
	return isSudden
}

// IsHighFrequencyTransaction detects "velocity" fraud by checking time intervals 
// between the current transaction and previous ones.
func (v *TxValidator) IsHighFrequencyTransaction(ctx context.Context, allPayments []models.Transaction, current models.Transaction) bool {
	currentTime, err := time.Parse(time.RFC3339, current.Timestamp)
	if err != nil {
		v.mu.RLock()
		if v.log != nil {
			v.log.Errorw(errors.ErrTxTimeParse.Error(), "raw_timestamp", current.Timestamp, "err", err)
		}
		v.mu.RUnlock()
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
			v.mu.RLock()
			if v.log != nil {
				v.log.Warnw(errors.ErrHighFreqTx.Error(), 
					"diff_sec", diff.Seconds(), 
					"limit_sec", minInterval.Seconds(),
				)
			}
			v.mu.RUnlock()
			return true
		}
	}
	return false
}

// GetStats returns the number of checks performed and fraud alerts triggered.
func (v *TxValidator) GetStats() (uint64, uint64) {
	return atomic.LoadUint64(&v.totalChecks), atomic.LoadUint64(&v.fraudAlerts)
}

// SetLogger safely updates the logger instance at runtime.
func (v *TxValidator) SetLogger(newLogger logPkg.Logger) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.log = newLogger
}