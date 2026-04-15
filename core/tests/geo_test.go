// Package tests implements validation suites for geographical scoring logic.
package tests

import (
	"context"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"

	"github.com/stretchr/testify/assert"
)

// TestLocationCalculators_AllCases provides comprehensive coverage for Country and City 
// verification. It ensures the engine correctly handles regional settings, string 
// normalization (trimming/casing), and weight distribution for both F2P and P2W profiles.
func TestLocationCalculators_AllCases(t *testing.T) {
	// Initialize calculators without logging for isolated logic testing.
	countryCalc := logic.RegCountryCalculator{}
	cityCalc := logic.RegCityCalculator{}

	type testCase struct {
		name      string
		userVal   string
		dbVal     string
		isDonator bool
		expected  float64
	}

	// Iterate through both user status profiles to verify weight scaling.
	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		// --- Country Validation Tests ---
		t.Run(prefix+"_Country", func(t *testing.T) {
			cases := []testCase{
				{"Match", "Moldova", "Moldova", isDonator, weights.RegCountry},
				{"CaseInsensitive", "moldova", "MOLDOVA", isDonator, weights.RegCountry},
				{"Mismatch", "Moldova", "Romania", isDonator, 0.0},
				{"EmptyUser", "", "Moldova", isDonator, 0.0},
				{"EmptyDB", "Moldova", "", isDonator, 0.0},
			}

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					user := models.UserData{UserClaim: models.UserClaim{RegCountry: c.userVal}}
					db := models.DBRecord{RegCountry: c.dbVal}

					// Using context.Background() to satisfy the interface.
					result := countryCalc.Calculate(context.Background(), user, db, weights)

					assert.InDelta(t, c.expected, result.Result, 0.001, "Value mismatch in %s", c.name)

					if c.expected > 0 {
						assert.Equal(t, "match", result.Status)
					} else if c.userVal == "" || c.dbVal == "" {
						assert.Equal(t, "no_data", result.Status)
					} else {
						assert.Equal(t, "no_match", result.Status)
					}
				})
			}
		})

		// --- City Validation Tests ---
		t.Run(prefix+"_City", func(t *testing.T) {
			cases := []testCase{
				{"Match", "Chisinau", "Chisinau", isDonator, weights.RegCity},
				{"WithSpaces", " Chisinau ", "Chisinau", isDonator, weights.RegCity}, // Verification of string trimming
				{"Mismatch", "Chisinau", "Balti", isDonator, 0.0},
				{"EmptyBoth", "", "", isDonator, 0.0},
			}

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					user := models.UserData{UserClaim: models.UserClaim{RegCity: c.userVal}}
					db := models.DBRecord{RegCity: c.dbVal}

					// Using context.Background() to satisfy the interface.
					result := cityCalc.Calculate(context.Background(), user, db, weights)

					assert.InDelta(t, c.expected, result.Result, 0.001, "Value mismatch in %s", c.name)

					if c.expected > 0 {
						assert.Equal(t, "match", result.Status)
					} else if c.userVal == "" || c.dbVal == "" {
						assert.Equal(t, "no_data", result.Status)
					} else {
						assert.Equal(t, "no_match", result.Status)
					}
				})
			}
		})
	}
}