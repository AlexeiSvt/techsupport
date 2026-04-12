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
	os.Setenv("API_IP_INFO_KEY", "test_key")
	defer os.Unsetenv("API_IP_INFO_KEY")

	now := time.Now()

	tests := []struct {
		name     string
		input    models.InputData
		expected float64
	}{
		{
			name: "01. Full Match (Local Trust IP)",
			input: models.InputData{
				UserData: models.UserData{
					IP: "127.0.0.1",
					UserClaim: models.UserClaim{
						AccTag:      "MR_ZIDDER",
						RegCountry:  "RU",
						RegCity:     "MSK",
						FirstEmail:  "test@mail.com",
						Phone:       "79991234567",
						FirstDevice: "PC",
						Devices:     []string{"PC"},
						RegDate:     now,
					},
				},
				DBRecord: models.DBRecord{
					RegCountry:  "RU",
					RegCity:     "MSK",
					FirstEmail:  "test@mail.com",
					Phone:       "79991234567",
					FirstDevice: "PC",
					Devices:     []string{"PC"},
					RegDate:     now,
					IsDonator:   false,
				},
			},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("\nCritical Error: %s\n Caused by: %v\n\nStack Trace:\n%s", tt.name, r, debug.Stack())
				}
			}()

			got := engine.CalculateFinalScore(tt.input)

			cleanPerc := strings.TrimSuffix(got.FinaPercentage, "%")
			gotValue, err := strconv.ParseFloat(cleanPerc, 64)
			if err != nil {
				t.Fatalf("Failed to parse percentage string: %s", got.FinaPercentage)
			}

			if math.Abs(gotValue-tt.expected) > 0.01 {
				t.Errorf("\nCalculation Error [%s]:\nExpected: %.2f\nGot (String): %s\nGot (Parsed): %.2f",
					tt.name, tt.expected, got.FinaPercentage, gotValue)
			}

			if len(got.Details) == 0 {
				t.Errorf("Error [%s]: Details slice is empty, engine didn't run calculators", tt.name)
			}
		})
	}
}
