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

// Фиксируем сид один раз для всего пакета, чтобы не спамить в цикле
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomDBRecord(i int) models.DBRecord {
	firstWindow := []models.Session{}
	lastWindow := []models.Session{}
	allPayments := []models.Transaction{}

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

	for j := 0; j < numPayments; j++ {
		allPayments = append(allPayments, models.Transaction{
			Amount:    float64(r.Intn(20)),
			Timestamp: time.Now().Add(-time.Duration(r.Intn(240)) * time.Hour).Format(time.RFC3339),
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
	return models.UserClaim{
		AccTag: "test_user",
		FirstTransaction: models.Transaction{
			City:      fmt.Sprintf("City-%d", r.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", r.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", r.Intn(50)),
			IP:        fmt.Sprintf("192.168.%d.%d", r.Intn(255), r.Intn(255)),
			Timestamp: time.Now().Format(time.RFC3339),
			Amount:    float64(r.Intn(100)),
		},
	}
}

func Test1000Transactions(t *testing.T) {
	calc := transactions.FirstTransactionScoreCalculator{}

	for i := 0; i < 1000; i++ {
		idx := i
		t.Run(fmt.Sprintf("tx_%d", idx), func(t *testing.T) {
			db := randomDBRecord(idx)
			claim := randomUserClaim()
			// ВАЖНО: Убедись, что GetWeights возвращает models.Weights
			w := logic.GetWeights(db.IsDonator)

			user := models.UserData{
				UserClaim: claim,
			}

			res := calc.Calculate(user, db, w)

			// ГАРАНТИРОВАННЫЙ ЛОГ: Никаких шансов для ошибки типов
			msg := fmt.Sprintf("[%v] Score: %.2f | Status: %v", res.Code, res.Result, res.Status)
			t.Log(msg)

			if res.Result > w.FirstTransaction {
				t.Errorf("Score %.2f > Weight %.2f", res.Result, w.FirstTransaction)
			}

			// Проверка на совпадение
			if len(db.UserHistory.FirstWindow) > 0 {
				user.UserClaim.FirstTransaction.City = db.UserHistory.FirstWindow[0].City
				user.UserClaim.FirstTransaction.Country = db.UserHistory.FirstWindow[0].Country
				user.UserClaim.FirstTransaction.DeviceID = db.UserHistory.FirstWindow[0].DeviceID
				
				matchRes := calc.Calculate(user, db, w)
				if matchRes.Result == 0 {
					t.Errorf("Expected positive score for matching data, got 0")
				}
			}
		})
	}
}