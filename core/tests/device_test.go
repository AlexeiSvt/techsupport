// Package tests contains logical verification for the scoring engine.
package tests

import (
	"context"
	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFirstDeviceCalculator_AllCases validates the hardware fingerprint matching logic.
// It ensures that original registration devices are correctly identified regardless
// of string casing, and that donor status weighting is properly applied.
func TestFirstDeviceCalculator_AllCases(t *testing.T) {
	// Initialize calculator. Logging is disabled for unit testing.
	calc := logic.FirstDeviceCalculator{}

	// Pre-load weight configurations for both user categories.
	wF2P := logic.GetWeights(false)
	wP2W := logic.GetWeights(true)

	// Define test matrix covering identity, mismatches, and empty data scenarios.
	cases := []struct {
		name      string
		userDev   string
		dbDev     string
		isDonator bool
		expected  float64
	}{
		// Standard Free-to-Play (F2P) cases
		{"F2P Identical", "iPhone 12", "iPhone 12", false, wF2P.FirstDevice},
		{"F2P Different", "iPhone 12", "Samsung Galaxy", false, 0.0},
		{"F2P EmptyBoth", "", "", false, 0.0},
		{"F2P CaseInsensitive", "iphone 12", "IPHONE 12", false, wF2P.FirstDevice},
		{"F2P Pixel", "Google Pixel 6", "Google Pixel 6", false, wF2P.FirstDevice},
		{"F2P LongDeviceName", "Samsung Galaxy S21 Ultra Premium Edition", "Samsung Galaxy S21 Ultra Premium Edition", false, wF2P.FirstDevice},

		// Pay-to-Win (P2W/Solvent) cases - ensuring weights scale correctly
		{"P2W Identical", "iPhone 12", "iPhone 12", true, wP2W.FirstDevice},
		{"P2W Different", "iPhone 12", "Samsung Galaxy", true, 0.0},
		{"P2W EmptyBoth", "", "", true, 0.0},
		{"P2W CaseInsensitive", "iphone 12", "IPHONE 12", true, wP2W.FirstDevice},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Retrieve the appropriate weight set based on the user's status.
			weights := logic.GetWeights(c.isDonator)

			user := models.UserData{
				UserClaim: models.UserClaim{FirstDevice: c.userDev},
			}
			db := models.DBRecord{
				FirstDevice: c.dbDev,
			}

			// Execution with background context to satisfy the ScoreCalculator interface.
			result := calc.Calculate(context.Background(), user, db, weights)

			// Assertions
			assert.InDelta(t, c.expected, result.Result, 0.001, "Value mismatch in case: %s", c.name)

			// Status logic verification:
			// 1. If result > 0, it must be a 'match'.
			// 2. If strings are empty, it must be 'no_data'.
			// 3. Otherwise, it is a 'no_match'.
			if c.expected > 0 {
				assert.Equal(t, "match", result.Status, "Status should be 'match' for: %s", c.name)
			} else if c.userDev == "" || c.dbDev == "" {
				assert.Equal(t, "no_data", result.Status, "Status should be 'no_data' for: %s", c.name)
			} else {
				assert.Equal(t, "no_match", result.Status, "Status should be 'no_match' for: %s", c.name)
			}
		})
	}
}