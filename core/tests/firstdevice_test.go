package tests

import (
	"fmt"

	"techsupport/core/internal/scoring/logic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScoreForFirstDevice_SingleDevice(t *testing.T) {
	refDevice := "iPhone 14"

	cases := []struct {
		name      string
		userDev   string
		dbDev     string
		isDonator bool
		expected  float64
	}{}

	for _, isDonator := range []bool{false, true} {
		weight := logic.GetWeights(isDonator).FirstDevice

		for i := 0; i < 50; i++ {
			// Совпадение
			cases = append(cases, struct {
				name      string
				userDev   string
				dbDev     string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("Match_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userDev:   refDevice,
				dbDev:     refDevice,
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, struct {
				name      string
				userDev   string
				dbDev     string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userDev:   "",
				dbDev:     refDevice,
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, struct {
				name      string
				userDev   string
				dbDev     string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userDev:   refDevice,
				dbDev:     "",
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, struct {
				name      string
				userDev   string
				dbDev     string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("Mismatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
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
			result := logic.CalculateScoreForFirstDevice(c.userDev, c.dbDev, weights)
			assert.Equal(t, c.expected, result, "Test failed: %s", c.name)
		})
	}
}
