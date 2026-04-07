package tests

import (
	"fmt"
	"techsupport/internal/scoring/logic"

	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateScoreForFirstPhone(t *testing.T) {
	refPhone := "+1234567890"

	cases := []struct {
		name       string
		userPhone  string
		dbPhone    string
		isDonator  bool
		expected   float64
	}{}

	for _, isDonator := range []bool{false, true} {
		weight := logic.GetWeights(isDonator).Phone

		for i := range 25 {
			cases = append(cases, struct {
				name      string
				userPhone string
				dbPhone   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("Match_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userPhone: refPhone,
				dbPhone:   refPhone,
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, struct {
				name      string
				userPhone string
				dbPhone   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("EmptyUser_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userPhone: "",
				dbPhone:   refPhone,
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, struct {
				name      string
				userPhone string
				dbPhone   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("EmptyDB_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userPhone: refPhone,
				dbPhone:   "",
				isDonator: isDonator,
				expected:  0,
			})

			cases = append(cases, struct {
				name      string
				userPhone string
				dbPhone   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("BothEmpty_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userPhone: "",
				dbPhone:   "",
				isDonator: isDonator,
				expected:  weight,
			})

			cases = append(cases, struct {
				name      string
				userPhone string
				dbPhone   string
				isDonator bool
				expected  float64
			}{
				name:      fmt.Sprintf("Mismatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
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
			result := logic.CalculateScoreForFirstPhone(c.userPhone, c.dbPhone, weights)
			assert.Equal(t, c.expected, result, "Test failed: %s", c.name)
		})
	}
}