// Package tests implements stress testing and property-based verification.
package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"techsupport/core/internal/logic"
	"techsupport/core/internal/logic/transactions"
	"techsupport/core/pkg/models"
)

// Use a fixed seed for reproducible test results.
var r = rand.New(rand.NewSource(42))

// randomDBRecord generates a mock database record with randomized session and payment history.
func randomDBRecord(i int) models.DBRecord {
	firstWindow := make([]models.Session, 0)
	lastWindow := make([]models.Session, 0)
	allPayments := make([]models.Transaction, 0)

	numFirst := 1 + r.Intn(3)
	numLast := 1 + r.Intn(3)
	numPayments := r.Intn(5)

	// Populate randomized session windows.
	for k := 0; k < numFirst; k++ {
		firstWindow = append(firstWindow, models.Session{
			City:      fmt.Sprintf("City-%d", r.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", r.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", r.Intn(50)),
			SessionIP: fmt.Sprintf("192.168.%d.%d", r.Intn(255), r.Intn(255)),
		})
	}

	for k := 0; k < numLast; k++ {
		lastWindow = append(lastWindow, models.Session{
			City:      fmt.Sprintf("City-%d", r.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", r.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", r.Intn(50)),
			SessionIP: fmt.Sprintf("192.168.%d.%d", r.Intn(255), r.Intn(255)),
		})
	}

	now := time.Now()

	// Populate randomized payment history.
	for j := 0; j < numPayments; j++ {
		allPayments = append(allPayments, models.Transaction{
			Amount:    float64(r.Intn(20)),
			Timestamp: now.Add(-time.Duration(r.Intn(240)) * time.Hour).Format(time.RFC3339),
		})
	}

	return models.DBRecord{
		IsDonator: i%2 == 0,
		UserHistory: models.UserHistory{
			FirstWindow: firstWindow,
			LastWindow:  lastWindow,
			AllPayments: allPayments,
		},
	}
}

// randomUserClaim generates a mock transaction claim from a user.
func randomUserClaim() models.UserClaim {
	now := time.Now()

	return models.UserClaim{
		AccTag: "test_user",
		FirstTransaction: models.Transaction{
			City:      fmt.Sprintf("City-%d", r.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", r.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", r.Intn(50)),
			IP:        fmt.Sprintf("192.168.%d.%d", r.Intn(255), r.Intn(255)),
			Timestamp: now.Format(time.RFC3339),
			Amount:    float64(r.Intn(100)),
		},
	}
}

// Test1000Transactions runs a high-concurrency stress test to ensure the 
// TransactionScoreCalculator is stable, thread-safe, and logically consistent 
// under diverse data inputs.
func Test1000Transactions(t *testing.T) {
	calc := transactions.FirstTransactionScoreCalculator{}

	for i := 0; i < 1000; i++ {
		// Capture variable for parallel execution.
		i := i

		t.Run(fmt.Sprintf("tx_%d", i), func(t *testing.T) {
			t.Parallel() // Enables concurrent execution of these 1000 cases.

			db := randomDBRecord(i)
			claim := randomUserClaim()
			w := logic.GetWeights(db.IsDonator)

			user := models.UserData{
				UserClaim: claim,
			}

			ctx := context.Background()
			res := calc.Calculate(ctx, user, db, w)

			// 1. Logic Check: Non-donators must always be skipped for transaction scoring.
			if !db.IsDonator {
				if res.Status != "skipped" {
					t.Fatalf("Constraint violation: expected status 'skipped' for non-donator, got %v", res.Status)
				}
				return
			}

			// 2. Bound Check: Result must never exceed the pre-defined weight.
			if res.Result > w.FirstTransaction {
				t.Fatalf("Score overflow: %.2f exceeds maximum weight %.2f", res.Result, w.FirstTransaction)
			}

			// 3. Consistency Check: If we force a match with the first session, score should be positive.
			if len(db.UserHistory.FirstWindow) > 0 {
				matchUser := user

				// Synthetically create a perfect match with the baseline window.
				matchUser.UserClaim.FirstTransaction.City = db.UserHistory.FirstWindow[0].City
				matchUser.UserClaim.FirstTransaction.Country = db.UserHistory.FirstWindow[0].Country
				matchUser.UserClaim.FirstTransaction.DeviceID = db.UserHistory.FirstWindow[0].DeviceID
				matchUser.UserClaim.FirstTransaction.IP = db.UserHistory.FirstWindow[0].SessionIP

				matchRes := calc.Calculate(ctx, matchUser, db, w)

				// Skip if blocked by other anomaly logic (e.g. amount too high).
				if matchRes.Status == "skipped" || matchRes.Status == "anomaly_block" {
					return
				}

				if matchRes.Result <= 0 {
					t.Fatalf("Logic failure: expected positive score for perfect historical match, got %.2f", matchRes.Result)
				}
			}
		})
	}
}