package tests

import (
	"fmt"

	"testing"

	"techsupport/internal/scoring/logic"

	"github.com/stretchr/testify/assert"
)

func TestCalculateScoreForRegCountryAndCity(t *testing.T) {
	refCountry := "Moldova"
	refCity := "Chisinau"

	cases := []struct {
		name       string
		userValue  string
		dbValue    string
		isDonator  bool
		expected   float64
		useCountry bool 
	}{}

	for _, isDonator := range []bool{false, true} {
		weightCountry := logic.GetWeights(isDonator).RegCountry
		weightCity := logic.GetWeights(isDonator).RegCity

		for i := range 25 { 

			cases = append(cases, struct {
				name       string
				userValue  string
				dbValue    string
				isDonator  bool
				expected   float64
				useCountry bool
			}{
				name:       fmt.Sprintf("CountryMatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userValue:  refCountry,
				dbValue:    refCountry,
				isDonator:  isDonator,
				expected:   weightCountry,
				useCountry: true,
			})

			cases = append(cases, struct {
				name       string
				userValue  string
				dbValue    string
				isDonator  bool
				expected   float64
				useCountry bool
			}{
				name:       fmt.Sprintf("CountryMismatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userValue:  refCountry,
				dbValue:    fmt.Sprintf("OtherCountry%d", i),
				isDonator:  isDonator,
				expected:   0,
				useCountry: true,
			})

			cases = append(cases, struct {
				name       string
				userValue  string
				dbValue    string
				isDonator  bool
				expected   float64
				useCountry bool
			}{
				name:       fmt.Sprintf("CityMatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userValue:  refCity,
				dbValue:    refCity,
				isDonator:  isDonator,
				expected:   weightCity,
				useCountry: false,
			})

			cases = append(cases, struct {
				name       string
				userValue  string
				dbValue    string
				isDonator  bool
				expected   float64
				useCountry bool
			}{
				name:       fmt.Sprintf("CityMismatch_%d_%s", i, map[bool]string{true: "P2W", false: "F2P"}[isDonator]),
				userValue:  refCity,
				dbValue:    fmt.Sprintf("OtherCity%d", i),
				isDonator:  isDonator,
				expected:   0,
				useCountry: false,
			})
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			weights := logic.GetWeights(c.isDonator)
			var result float64
			if c.useCountry {
				result = logic.CalculateScoreForRegCountry(c.userValue, c.dbValue, weights)
			} else {
				result = logic.CalculateScoreForRegCity(c.userValue, c.dbValue, weights)
			}
			assert.Equal(t, c.expected, result, "Test failed: %s", c.name)
		})
	}
}
