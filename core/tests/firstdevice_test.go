package tests

import (
	"fmt"
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFirstDeviceCalculator_SingleDevice(t *testing.T) {
	refDevice := "iPhone 14"
	calc := logic.FirstDeviceCalculator{}

	type testCase struct {
		name      string
		userDev   string
		dbDev     string
		isDonator bool
		expected  float64
	}

	var cases []testCase

	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		weight := weights.FirstDevice
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		for i := range 10 {
			cases = append(cases, testCase{
				name:      fmt.Sprintf("Match_%d_%s", i, prefix),
				userDev:   refDevice,
				dbDev:     refDevice,
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, prefix),
				userDev:   "",
				dbDev:     refDevice,
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, testCase{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, prefix),
				userDev:   refDevice,
				dbDev:     "",
				isDonator: isDonator,
				expected:  0,
			})

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
			weights := logic.GetWeights(c.isDonator)

			user := models.UserData{
				UserClaim: models.UserClaim{FirstDevice: c.userDev},
			}
			db := models.DBRecord{
				FirstDevice: c.dbDev,
			}

			result := calc.Calculate(user, db, weights)

			assert.InDelta(t, c.expected, result.Result, 0.001, "Test failed: %s", c.name)

			if c.expected > 0 {
				assert.Equal(t, "match", result.Status)
			} else if c.userDev == "" || c.dbDev == "" {
				assert.Equal(t, "no_data", result.Status)
			} else {
				assert.Equal(t, "no_match", result.Status)
			}
		})
	}
}