package tests

import (
	"fmt"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFirstEmailCalculator_AllCases(t *testing.T) {
	refEmail := "user@example.com"
	calc := logic.FirstEmailCalculator{}

	type testCase struct {
		name      string
		userEmail string
		dbEmail   string
		isDonator bool
		expected  float64
	}

	var cases []testCase

	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		weight := weights.FirstEmail
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		for i := range 5 { 
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Match_%d_%s", i, prefix),
				userEmail: refEmail,
				dbEmail:   refEmail,
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, prefix),
				userEmail: "",
				dbEmail:   refEmail,
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, prefix),
				userEmail: refEmail,
				dbEmail:   "",
				isDonator: isDonator,
				expected:  0,
			})

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

			result := calc.Calculate(user, db, weights)

			assert.InDelta(t, c.expected, result.Result, 0.001, "Test failed: %s", c.name)

			if c.expected > 0 {
				assert.Equal(t, "match", result.Status)
			} else if c.userEmail == "" || c.dbEmail == "" {
				assert.Equal(t, "no_data", result.Status)
			}
		})
	}
}