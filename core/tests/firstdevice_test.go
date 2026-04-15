// Package tests implements validation suites for hardware fingerprinting logic.
package tests

import (
	"context"
	"fmt"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"
	"github.com/stretchr/testify/assert"
)

// TestFirstDeviceCalculator_SingleDevice executes an exhaustive suite of tests
// for the primary device verification logic. It ensures that both free-to-play (F2P)
// and donor (P2W) weighting profiles handle matches, mismatches, and missing data consistently.
func TestFirstDeviceCalculator_SingleDevice(t *testing.T) {
	refDevice := "iPhone 14"
	// Initialize calculator. Logging is disabled for low-level unit verification.
	calc := logic.FirstDeviceCalculator{}

	type testCase struct {
		name      string
		userDev   string
		dbDev     string
		isDonator bool
		expected  float64
	}

	var cases []testCase

	// Iterate through both user status profiles to verify weight scaling.
	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		weight := weights.FirstDevice
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		// Perform multiple iterations to ensure logic consistency and stress-test the data structures.
		for i := 0; i < 10; i++ {
			// Scenario: Exact hardware fingerprint match.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Match_%d_%s", i, prefix),
				userDev:   refDevice,
				dbDev:     refDevice,
				isDonator: isDonator,
				expected:  weight,
			})

			// Scenario: User failed to provide a device claim.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, prefix),
				userDev:   "",
				dbDev:     refDevice,
				isDonator: isDonator,
				expected:  0,
			})

			// Scenario: Historical database record is missing device info.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, prefix),
				userDev:   refDevice,
				dbDev:     "",
				isDonator: isDonator,
				expected:  0,
			})

			// Scenario: Explicit mismatch between current device and registration device.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Mismatch_%d_%s", i, prefix),
				userDev:   refDevice,
				dbDev:     "Samsung Galaxy S23",
				isDonator: isDonator,
				expected:  0,
			})
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Retrieve the dynamic weight profile for the current test case.
			weights := logic.GetWeights(c.isDonator)

			user := models.UserData{
				UserClaim: models.UserClaim{FirstDevice: c.userDev},
			}
			db := models.DBRecord{
				FirstDevice: c.dbDev,
			}

			// Execution using context.Background() to satisfy the Engine's calculator interface.
			result := calc.Calculate(context.Background(), user, db, weights)

			// Floating point comparison with delta for reliability.
			assert.InDelta(t, c.expected, result.Result, 0.001, "Scoring value mismatch: %s", c.name)

			// Categorical status verification:
			// - Match: successful verification.
			// - No Data: one or both inputs were empty strings.
			// - No Match: explicit data conflict.
			if c.expected > 0 {
				assert.Equal(t, "match", result.Status, "Expected status 'match' for case: %s", c.name)
			} else if c.userDev == "" || c.dbDev == "" {
				assert.Equal(t, "no_data", result.Status, "Expected status 'no_data' for case: %s", c.name)
			} else {
				assert.Equal(t, "no_match", result.Status, "Expected status 'no_match' for case: %s", c.name)
			}
		})
	}
}