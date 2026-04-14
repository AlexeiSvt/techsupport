package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"techsupport/core/internal/logic"
	"techsupport/core/internal/logic/transactions"
	"techsupport/core/internal/models"
)

var r = rand.New(rand.NewSource(42))

func randomDBRecord(i int) models.DBRecord {
	firstWindow := make([]models.Session, 0)
	lastWindow := make([]models.Session, 0)
	allPayments := make([]models.Transaction, 0)

	numFirst := 1 + r.Intn(3)
	numLast := 1 + r.Intn(3)
	numPayments := r.Intn(5)

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

func Test1000Transactions(t *testing.T) {
	calc := transactions.FirstTransactionScoreCalculator{Log: nil}

	for i := 0; i < 1000; i++ {
		i := i

		t.Run(fmt.Sprintf("tx_%d", i), func(t *testing.T) {
			t.Parallel()

			db := randomDBRecord(i)
			claim := randomUserClaim()
			w := logic.GetWeights(db.IsDonator)

			user := models.UserData{
				UserClaim: claim,
			}

			res := calc.Calculate(user, db, w)

			if !db.IsDonator {
				if res.Status != "skipped" {
					t.Fatalf("expected status skipped for non-donator, got %v", res.Status)
				}
				return
			}

			if res.Result > w.FirstTransaction {
				t.Fatalf("score %.2f exceeds weight %.2f", res.Result, w.FirstTransaction)
			}

			if len(db.UserHistory.FirstWindow) > 0 {
				matchUser := user

				matchUser.UserClaim.FirstTransaction.City = db.UserHistory.FirstWindow[0].City
				matchUser.UserClaim.FirstTransaction.Country = db.UserHistory.FirstWindow[0].Country
				matchUser.UserClaim.FirstTransaction.DeviceID = db.UserHistory.FirstWindow[0].DeviceID
				matchUser.UserClaim.FirstTransaction.IP = db.UserHistory.FirstWindow[0].SessionIP

				matchRes := calc.Calculate(matchUser, db, w)

				if matchRes.Status == "skipped" || matchRes.Status == "anomaly_block" {
					return
				}

				if matchRes.Result <= 0 {
					t.Fatalf("expected positive score for match, got %.2f", matchRes.Result)
				}
			}
		})
	}
}