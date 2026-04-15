// Package tests ensures the reliability of the scoring logic through rigorous unit testing.
package tests

import (
	"context"
	"fmt"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"
	"github.com/stretchr/testify/assert"
)

// TestFirstEmailCalculator_AllCases validates the email matching logic across 
// multiple scenarios, ensuring that both free-to-play and donor weightings 
// are correctly applied and that empty or mismatched data triggers appropriate statuses.
func TestFirstEmailCalculator_AllCases(t *testing.T) {
	refEmail := "user@example.com"
	// Initialize calculator. Logging is disabled for isolated logic tests.
	calc := logic.FirstEmailCalculator{}

	type testCase struct {
		name      string
		userEmail string
		dbEmail   string
		isDonator bool
		expected  float64
	}

	var cases []testCase

	// Generate a matrix of test cases for both user categories.
	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		weight := weights.FirstEmail
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		// Using a loop to generate multiple variations of each scenario type.
		for i := range 5 { 
			// 1. Success case: Emails match perfectly.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Match_%d_%s", i, prefix),
				userEmail: refEmail,
				dbEmail:   refEmail,
				isDonator: isDonator,
				expected:  weight,
			})

			// 2. Data missing case: User provided an empty claim.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, prefix),
				userEmail: "",
				dbEmail:   refEmail,
				isDonator: isDonator,
				expected:  0,
			})

			// 3. Data missing case: Database record is empty.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, prefix),
				userEmail: refEmail,
				dbEmail:   "",
				isDonator: isDonator,
				expected:  0,
			})

			// 4. Mismatch case: Provided email does not match the historical record.
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Mismatch_%d_%s", i, prefix),
				userEmail: refEmail,
				dbEmail:   fmt.Sprintf("other%d@example.com", i),
				isDonator: isDonator,
				expected:  0,
			})
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			weights := logic.GetWeights(c.isDonator)
			
			user := models.UserData{
				UserClaim: models.UserClaim{FirstEmail: c.userEmail},
			}
			db := models.DBRecord{
				FirstEmail: c.dbEmail,
			}

			// Added context.Background() to satisfy the interface requirements.
			result := calc.Calculate(context.Background(), user, db, weights)

			// Verify numerical precision with InDelta.
			assert.InDelta(t, c.expected, result.Result, 0.001, "Score mismatch in: %s", c.name)

			// Status logic verification based on business requirements:
			// - Points earned: "match"
			// - One or both fields empty: "no_data"
			// - Explicit mismatch: "no_match"
			if c.expected > 0 {
				assert.Equal(t, "match", result.Status, "Expected status 'match' for: %s", c.name)
			} else if c.userEmail == "" || c.dbEmail == "" {
				assert.Equal(t, "no_data", result.Status, "Expected status 'no_data' for: %s", c.name)
			} else {
				assert.Equal(t, "no_match", result.Status, "Expected status 'no_match' for: %s", c.name)
			}
		})
	}
}