package tests

import (
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"techsupport/core/internal/engine"
	"techsupport/core/internal/models"
)

func TestCalculateFinalScore_Production(t *testing.T) {
	apiKey := os.Getenv("API_IP_INFO_KEY")
	if apiKey == "" {
		t.Skip("Пропуск: API_IP_INFO_KEY is not set. Set it to run production tests.")
    }

	now := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name     string
		input    models.InputData
		expected float64 // Ожидаемый итоговый процент (0.00 - 100.00)
	}{
		{
			name: "01. Full Match (Localhost - Bogon IP)",
			input: models.InputData{
				UserData: models.UserData{
					IP: "127.0.0.1",
					UserClaim: models.UserClaim{
						AccTag:      "ALEXEY_DEV",
						RegCountry:  "MD",
						RegCity:     "Chisinau",
						FirstEmail:  "alexey@test.com",
						Phone:       "37360000000",
						FirstDevice: "PC",
						Devices:     []string{"PC"},
						RegDate:     now,
					},
				},
				DBRecord: models.DBRecord{
					RegCountry:  "MD",
					RegCity:     "Chisinau",
					FirstEmail:  "alexey@test.com",
					Phone:       "37360000000",
					FirstDevice: "PC",
					Devices:     []string{"PC"},
					RegDate:     now,
				},
			},
			// Для реального API адрес 127.0.0.1 — это Bogon. 
			// Штраф будет 100%, итоговый балл 0.
			expected: 0.00,
		},
		{
			name: "02. Partial Match (Google Public DNS - Datacenter)",
			input: models.InputData{
				UserData: models.UserData{
					IP: "8.8.8.8",
					UserClaim: models.UserClaim{
						AccTag:      "ALEXEY_DEV",
						RegCountry:  "US",
						RegCity:     "Mountain View",
						FirstEmail:  "alexey@test.com",
						Phone:       "37360000000",
						FirstDevice: "PC",
						Devices:     []string{"PC"},
						RegDate:     now,
					},
				},
				DBRecord: models.DBRecord{
					RegCountry:  "US",
					RegCity:     "Mountain View",
					FirstEmail:  "alexey@test.com",
					Phone:       "37360000000",
					FirstDevice: "PC",
					Devices:     []string{"PC"},
					RegDate:     now,
				},
			},
			expected: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Защита от паники (Recover)
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("\nКритический сбой в [%s]: %v\nСтек:\n%s", tt.name, r, debug.Stack())
				}
			}()

			got := engine.CalculateFinalScore(tt.input)

			cleanPerc := strings.TrimSuffix(got.FinaPercentage, "%")
			gotValue, err := strconv.ParseFloat(cleanPerc, 64)
			if err != nil {
				t.Fatalf("Ошибка парсинга процента: %s", got.FinaPercentage)
			}

			if math.Abs(gotValue-tt.expected) > 0.01 {
				t.Errorf("\nОшибка расчета [%s]:\nОжидали: %.2f\nПолучили: %s\nОператор IP: %s",
					tt.name, tt.expected, got.FinaPercentage, got.Details[0].MetricName) // Предполагаем, что ASN в деталях
			}

			if len(got.Details) == 0 {
				t.Errorf("Ошибка [%s]: Список Details пуст, расчеты не проводились", tt.name)
			}
		})
	}
}