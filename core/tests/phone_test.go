// Package tests implements logical verification for identity marker calculators.
package tests

import (
	"context"
	"fmt"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"
	"github.com/stretchr/testify/assert"
)

// TestFirstPhoneCalculator_AllCases ensures that phone number verification
// correctly applies weights for both standard (F2P) and donor (P2W) profiles.
// It validates exact matches, empty inputs, and explicit mismatches.
func TestFirstPhoneCalculator_AllCases(t *testing.T) {
	refPhone := "+1234567890"
	// Initialize calculator. Logging is disabled for core logic unit testing.
	calc := logic.FirstPhoneCalculator{}

	type testCase struct {
		name      string
		userPhone string
		dbPhone   string
		isDonator bool
		expected  float64
	}

	var cases []testCase

	// Iterate through user categories to ensure weight sets are correctly retrieved.
	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		weight := weights.Phone 
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		// Generate variations for comprehensive coverage.
		for i := 0; i < 5; i++ {
			// Scenario: Successful identity verification via phone match.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Match_%d_%s", i, prefix),
				userPhone: refPhone,
				dbPhone:   refPhone,
				isDonator: isDonator,
				expected:  weight,
			})

			// Scenario: Missing claim from the user's side.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, prefix),
				userPhone: "",
				dbPhone:   refPhone,
				isDonator: isDonator,
				expected:  0,
			})

			// Scenario: Phone number mismatch (potential account takeover or wrong data).
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Mismatch_%d_%s", i, prefix),
				userPhone: refPhone,
				dbPhone:   fmt.Sprintf("+1000000%d", i),
				isDonator: isDonator,
				expected:  0,
			})
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			weights := logic.GetWeights(c.isDonator)
			
			user := models.UserData{
				UserClaim: models.UserClaim{Phone: c.userPhone},
			}
			db := models.DBRecord{
				Phone: c.dbPhone,
			}

			// Execution with background context to satisfy the calculator interface.
			// Note: If your current logic doesn't take context yet, remove context.Background().
			result := calc.Calculate(context.Background(), user, db, weights)

			// Validate numerical result with float delta.
			assert.InDelta(t, c.expected, result.Result, 0.001, "Scoring mismatch in: %s", c.name)

			// Status-based verification logic.
			if c.expected > 0 {
				assert.Equal(t, "match", result.Status, "Status should be 'match' for: %s", c.name)
			} else if c.userPhone == "" || c.dbPhone == "" {
				assert.Equal(t, "no_data", result.Status, "Status should be 'no_data' for: %s", c.name)
			} else {
				assert.Equal(t, "no_match", result.Status, "Status should be 'no_match' for: %s", c.name)
			}
		})
	}
}