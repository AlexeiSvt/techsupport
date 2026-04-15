// Package tests contains integration and unit tests for the scoring logic.
package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"

	"github.com/stretchr/testify/assert"
)

// TestCalculateScoreForCreationAge_RegDate_BoundaryCases performs exhaustive testing
// of the account age validation logic. It covers various offsets in hours and months
// to ensure that the "match", "partial", and "no_match" thresholds are razor-sharp.
func TestCalculateScoreForCreationAge_RegDate_BoundaryCases(t *testing.T) {
	// Fixed reference date for consistent testing.
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// RegDateCalculator is stateless, so we can initialize it without a logger for tests.
	calc := logic.RegDateCalculator{} 

	type testCase struct {
		name           string
		userDate       time.Time
		dbDate         time.Time
		isDonator      bool
		expectedResult float64
		expectedStatus string
	}

	var cases []testCase

	// Test for both Donator (Solvent) and Non-Donator (F2P) users as they might have different weights.
	isDonatorOptions := []bool{false, true}
	for _, isDonator := range isDonatorOptions {
		weights := logic.GetWeights(isDonator)
		weight := weights.RegDate

		// Generates a matrix of time offsets:
		// i: hour jitter (-5 to +5) to ensure time-of-day doesn't break the month logic.
		// j: month difference (0 to 3) to test threshold boundaries.
		for i := -5; i <= 5; i++ {
			for j := 0; j <= 3; j++ {
				expectedRes := 0.0
				expectedStatus := "no_match"
				diffMonths := float64(j)
				
				// Threshold logic verification:
				// - Within 2 months: Full points (match)
				// - 3 to 4 months: Half points (partial)
				// - Beyond 4 months: Zero points (no_match)
				if diffMonths <= 2 {
					expectedRes = weight 
					expectedStatus = "match"
				} else if diffMonths <= 4 {
					expectedRes = weight * 0.5 
					expectedStatus = "partial"
				}

				userTime := baseDate.Add(time.Duration(i) * time.Hour).AddDate(0, j, 0)
				
				prefix := "F2P"
				if isDonator { prefix = "Solvent" }

				cases = append(cases, testCase{
					name:           fmt.Sprintf("%s_hours=%d_months=%d", prefix, i, j),
					userDate:       userTime,
					dbDate:         baseDate,
					isDonator:      isDonator,
					expectedResult: expectedRes,
					expectedStatus: expectedStatus,
				})
			}
		}
	}

	// Execution of the generated test suite.
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			weights := logic.GetWeights(c.isDonator)
			
			user := models.UserData{
				UserClaim: models.UserClaim{RegDate: c.userDate},
			}
			db := models.DBRecord{
				RegDate: c.dbDate,
			}

			// Context is omitted here if the calculator doesn't use it for RegDate,
			// or passed as context.Background() if the interface requires it.
			result := calc.Calculate(context.Background(), user, db, weights)

			// InDelta is used for float comparison to avoid precision issues.
			assert.InDelta(t, c.expectedResult, result.Result, 0.001, "Result mismatch in test case: %s", c.name)
			
			// Exact string match for status.
			assert.Equal(t, c.expectedStatus, result.Status, "Status mismatch in test case: %s", c.name)
		})
	}
}