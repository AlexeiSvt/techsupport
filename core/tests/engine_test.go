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
	"techsupport/core/internal/ipchecker"
	"techsupport/core/pkg/models"
)

func TestCalculateFinalScore_Production(t *testing.T) {
	apiKey := os.Getenv("API_IP_INFO_KEY")
	if apiKey == "" {
		t.Skip("SKIP: API_IP_INFO_KEY is not set.")
	}

	now := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name     string
		input    models.InputData
		expected float64
	}{
		{
			name: "02. Google DNS (Datacenter)",
			input: models.InputData{
				UserData: models.UserData{
					IPInfo: &ipchecker.IpApiResponse{
						IP:           "8.8.8.8",
						TrustScore:   90,
						IsDatacenter: true,
						IsVPN:        false,
						IsMobile:     false,
						ASN: ipchecker.ASNInfo{
							Number: "AS15169",
							Org: "Google LLC",
						},
					},
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
			expected: 0.00,
		},
	} 

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic: %v\n%s", r, debug.Stack())
				}
			}()

			ctx := context.Background()
			got := engine.CalculateFinalScore(ctx, nil, tt.input)

			clean := strings.TrimSuffix(got.FinaPercentage, "%")
			val, err := strconv.ParseFloat(clean, 64)
			if err != nil {
				t.Fatalf("parse error: %s", got.FinaPercentage)
			}

			if math.Abs(val-tt.expected) > 0.01 {
				t.Errorf("expected %.2f, got %s", tt.expected, got.FinaPercentage)
			}

			if len(got.Details) == 0 {
				t.Errorf("empty details")
			}
		})
	}
}