// Package tests implements high-level integration testing for the scoring engine.
package tests

import (
	"context"
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"techsupport/core/internal/engine"
	"techsupport/core/pkg/models"
)

// TestCalculateFinalScore_Production simulates real-world scenarios to verify 
// the end-to-end scoring pipeline, including IP intelligence and penalty application.
func TestCalculateFinalScore_Production(t *testing.T) {
	// Ensure we have the API key for external IP reputation checks.
	apiKey := os.Getenv("API_IP_INFO_KEY")
	if apiKey == "" {
		t.Skip("SKIP: API_IP_INFO_KEY is not set. This test requires live API access.")
	}

	// Reference time for consistency (yesterday).
	now := time.Now().Add(-24 * time.Hour)

	// Define production-like test scenarios.
	tests := []struct {
		name     string
		input    models.InputData
		expected float64 // Expected final percentage value.
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
						FirstEmail:  "test@mail.com",
						Phone:       "37360000000",
						FirstDevice: "PC",
						Devices:     []string{"PC"},
						RegDate:     now,
					},
				},
				DBRecord: models.DBRecord{
					IsDonator:   false,
					RegCountry:  "MD",
					RegCity:     "Chisinau",
					FirstEmail:  "test@mail.com",
					Phone:       "37360000000",
					FirstDevice: "PC",
					Devices:     []string{"PC"},
					RegDate:     now,
				},
			},
			expected: 0.00, // Adjust this based on your specific logic for bogon IPs.
		},
		{
			name: "02. Google DNS (Datacenter)",
			input: models.InputData{
				UserData: models.UserData{
					IP: "8.8.8.8",
					UserClaim: models.UserClaim{
						AccTag:      "ALEXEY_DEV",
						RegCountry:  "US",
						RegCity:     "Mountain View",
						FirstEmail:  "test@mail.com",
						Phone:       "37360000000",
						FirstDevice: "PC",
						Devices:     []string{"PC"},
						RegDate:     now,
					},
				},
				DBRecord: models.DBRecord{
					IsDonator:   false,
					RegCountry:  "US",
					RegCity:     "Mountain View",
					FirstEmail:  "test@mail.com",
					Phone:       "37360000000",
					FirstDevice: "PC",
					Devices:     []string{"PC"},
					RegDate:     now,
				},
			},
			expected: 0.00, // Adjust based on Datacenter penalty logic.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Graceful recovery from unexpected panics during complex calculations.
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("\nCRITICAL FAILURE in [%s]: %v\nStack Trace:\n%s", tt.name, r, debug.Stack())
				}
			}()

			// Using Background context and nil for logger (or a mock logger if available).
			ctx := context.Background()
			// FIX: Passing context and nil logger to match the updated function signature.
			got := engine.CalculateFinalScore(ctx, nil, tt.input)

			// Parse "XX.XX%" string back to float for numerical comparison.
			cleanPerc := strings.TrimSuffix(got.FinaPercentage, "%")
			gotValue, err := strconv.ParseFloat(cleanPerc, 64)
			if err != nil {
				t.Fatalf("Failed to parse final percentage: %s", got.FinaPercentage)
			}

			// Validate with a small delta to account for floating point precision.
			if math.Abs(gotValue-tt.expected) > 0.01 {
				t.Errorf("\nCalculation Logic Error [%s]:\nExpected: %.2f\nGot: %s\nDetails Trace: %+v",
					tt.name, tt.expected, got.FinaPercentage, got.Details)
			}

			// Ensure the audit trail is being populated.
			if len(got.Details) == 0 {
				t.Errorf("Audit Error [%s]: Calculation details slice is empty", tt.name)
			}
		})
	}
}