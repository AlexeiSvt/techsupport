package tests

import (
	"fmt"
	"testing"
	"time"

	"techsupport/internal/scoring/logic"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDeltaOfCreationAge_RegDate_BoundaryCases(t *testing.T) {
	baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name      string
		userDate  time.Time
		dbDate    time.Time
		isDonator bool
		expected  float64
	}{}

	for i := -5; i <= 5; i++ {
		for j := 0; j <= 3; j++ { 
			expected := 0.0
			diffMonths := float64(j)
			if diffMonths <= 2 {
				expected = 12.5
			} else if diffMonths <= 4 {
				expected = 12.5 * 0.5 
			}

			userTime := baseDate.Add(time.Duration(i) * time.Hour).AddDate(0, j, 0)
			cases = append(cases, struct {
				name      string
				userDate  time.Time
				dbDate    time.Time
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("F2P_hours=%d_months=%d", i, j),
				userDate:  userTime,
				dbDate:    baseDate,
				isDonator: false,
				expected:  expected,
			})
		}
	}


	for i := -5; i <= 5; i++ {
		for j := 0; j <= 3; j++ {
			expected := 0.0
			diffMonths := float64(j)
			if diffMonths <= 2 {
				expected = 7.5 
			} else if diffMonths <= 4 {
				expected = 7.5 * 0.5
			}

			userTime := baseDate.Add(time.Duration(i) * time.Hour).AddDate(0, j, 0)
			cases = append(cases, struct {
				name      string
				userDate  time.Time
				dbDate    time.Time
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("P2W_hours=%d_months=%d", i, j),
				userDate:  userTime,
				dbDate:    baseDate,
				isDonator: true,
				expected:  expected,
			})
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			weights := logic.GetWeights(c.isDonator)
			result := logic.CalculateDeltaOfCreationAge(c.userDate, c.dbDate, weights)
			assert.Equal(t, c.expected, result, "Test failed: %s", c.name)
		})
	}
}