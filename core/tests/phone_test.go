package tests

import (
	"fmt"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFirstPhoneCalculator_AllCases(t *testing.T) {
	refPhone := "+1234567890"
	calc := logic.FirstPhoneCalculator{Log: nil}

	type testCase struct {
		name      string
		userPhone string
		dbPhone   string
		isDonator bool
		expected  float64
	}

	var cases []testCase

	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		weight := weights.Phone 
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		for i := range 5 {
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Match_%d_%s", i, prefix),
				userPhone: refPhone,
				dbPhone:   refPhone,
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, prefix),
				userPhone: "",
				dbPhone:   refPhone,
				isDonator: isDonator,
				expected:  0,
			})

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

			result := calc.Calculate(user, db, weights)

			assert.InDelta(t, c.expected, result.Result, 0.001, "Value mismatch in: %s", c.name)

			if c.expected > 0 {
				assert.Equal(t, "match", result.Status, "Should be match: %s", c.name)
			} else if c.userPhone == "" || c.dbPhone == "" {
				assert.Equal(t, "no_data", result.Status, "Should be no_data: %s", c.name)
			} else {
				assert.Equal(t, "no_match", result.Status, "Should be no_match: %s", c.name)
			}
		})
	}
}