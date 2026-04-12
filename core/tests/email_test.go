package tests

import (
	"fmt"
	"testing"

	"techsupport/core/internal/logic"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScoreForFirstEmail_SingleEmail(t *testing.T) {
	refEmail := "user@example.com"

	cases := []struct {
		name      string
		userEmail string
		dbEmail   string
		isDonator bool
		expected  float64
	}{}

	for _, isDonator := range []bool{false, true} {
		weight := logic.GetWeights(isDonator).FirstEmail
		for i := range 25 {
			cases = append(cases, struct {
				name      string
				userEmail string
				dbEmail   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("Match_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userEmail: refEmail,
				dbEmail:   refEmail,
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, struct {
				name      string
				userEmail string
				dbEmail   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userEmail: "",
				dbEmail:   refEmail,
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, struct {
				name      string
				userEmail string
				dbEmail   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userEmail: refEmail,
				dbEmail:   "",
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, struct {
				name      string
				userEmail string
				dbEmail   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("Mismatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
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
			result := logic.CalculateScoreForFirstEmail(c.userEmail, c.dbEmail, weights)
			assert.Equal(t, c.expected, result, "Test failed: %s", c.name)
		})
	}
}