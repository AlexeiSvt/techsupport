package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"techsupport/core/internal/models"
	"techsupport/core/internal/scoring/logic"
	"techsupport/core/internal/scoring/logic/transactions"
)

func randomDBRecord(i int) models.DBRecord {
	firstWindow := []models.Session{}
	lastWindow := []models.Session{}
	allPayments := []models.Transaction{}

	numFirst := 1 + rand.Intn(3)
	numLast := 1 + rand.Intn(3)
	numPayments := rand.Intn(5)

	for range numFirst {
		firstWindow = append(firstWindow, models.Session{
			City:      fmt.Sprintf("City-%d", rand.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", rand.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", rand.Intn(50)),
			SessionIP: fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
			ASN:       []string{"Rostelecom", "Beeline", "MTS", "Tele2"}[rand.Intn(4)],
		})
	}

	for range numLast {
		lastWindow = append(lastWindow, models.Session{
			City:      fmt.Sprintf("City-%d", rand.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", rand.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", rand.Intn(50)),
			SessionIP: fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
			ASN:       []string{"OVH", "DigitalOcean", "Hetzner", "AWS"}[rand.Intn(4)],
		})
	}

	for j := 0; j < numPayments; j++ {
		allPayments = append(allPayments, models.Transaction{
			City:      fmt.Sprintf("City-%d", rand.Intn(10)),
			Country:   fmt.Sprintf("Country-%d", rand.Intn(5)),
			DeviceID:  fmt.Sprintf("device-%d", rand.Intn(50)),
			IP:        fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
			Timestamp: time.Now().Add(-time.Duration(rand.Intn(240)) * time.Hour).Format(time.RFC3339),
			Amount:    float64(rand.Intn(20)), // донаты до 20$
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
	cities := []string{"City-0", "City-1", "City-2", "City-3", "City-4", "City-5", "City-6", "City-7", "City-8", "City-9"}
	countries := []string{"Country-0", "Country-1", "Country-2", "Country-3", "Country-4"}

	return models.UserClaim{
		FirstTransaction: models.Transaction{
			City:      cities[rand.Intn(len(cities))],
			Country:   countries[rand.Intn(len(countries))],
			DeviceID:  fmt.Sprintf("device-%d", rand.Intn(50)),
			IP:        fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
			ASN:       []string{"Rostelecom", "Beeline", "MTS", "Tele2", "OVH", "AWS"}[rand.Intn(6)],
			Timestamp: time.Now().Format(time.RFC3339),
			Amount:    float64(rand.Intn(50)),
		},
		Devices: []string{
			fmt.Sprintf("device-%d", rand.Intn(50)),
			fmt.Sprintf("device-%d", rand.Intn(50)),
		},
	}
}

func Test1000Transactions(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := range 1000 {
		t.Run(fmt.Sprintf("tx_%d", i), func(t *testing.T) {
			db := randomDBRecord(i)
			user := randomUserClaim()
			weights := logic.GetWeights(db.IsDonator)

			if weights.FirstTransaction <= 0 {
				t.Skip()
			}

			score := transactions.CalculateFirstTransactionScore(db, user, weights)
			t.Logf("tx_%d: score=%.3f, firstWeight=%.3f, amount=%.2f", i, score, weights.FirstTransaction, user.FirstTransaction.Amount)

			// Проверка диапазона
			if score < -500 || score > weights.FirstTransaction {
				t.Errorf("score out of expected range: %f", score)
			}


			user.FirstTransaction.City = db.UserHistory.LastWindow[0].City
			user.FirstTransaction.DeviceID = "unknown-device"
			partialScore := transactions.CalculateFirstTransactionScore(db, user, weights)
			if partialScore <= 0 {
				t.Logf("DEBUG: partial match <=0: score=%.3f", partialScore)
			}

			if isSuddenHighDonation(user.FirstTransaction, db.UserHistory) {
				t.Logf("DEBUG: sudden high donation, amount=%.2f, score=%.3f", user.FirstTransaction.Amount, score)
			}
		})
	}
}


func isSuddenHighDonation(tx models.Transaction, history models.UserHistory) bool {
	const suddenMultiplier = 5.0
	const firstDonationThreshold = 5.0

	if len(history.AllPayments) == 0 {
		return tx.Amount >= firstDonationThreshold
	}

	var total float64
	for _, p := range history.AllPayments {
		total += p.Amount
	}
	avg := total / float64(len(history.AllPayments))
	return tx.Amount > avg*suddenMultiplier
}
