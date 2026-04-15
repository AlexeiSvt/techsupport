package tests

import (
	"testing"

	"techsupport/core/internal/logic"
	"techsupport/core/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestLocationCalculators_AllCases(t *testing.T) {
	countryCalc := logic.RegCountryCalculator{Log: nil}
	cityCalc := logic.RegCityCalculator{Log: nil}

	type testCase struct {
		name      string
		userVal   string
		dbVal     string
		isDonator bool
		expected  float64
	}

	for _, isDonator := range []bool{false, true} {
		weights := logic.GetWeights(isDonator)
		prefix := map[bool]string{true: "P2W", false: "F2P"}[isDonator]

		// Тесты для Страны
		t.Run(prefix+"_Country", func(t *testing.T) {
			cases := []testCase{
				{"Match", "Moldova", "Moldova", isDonator, weights.RegCountry},
				{"CaseInsensitive", "moldova", "MOLDOVA", isDonator, weights.RegCountry},
				{"Mismatch", "Moldova", "Romania", isDonator, 0.0},
				{"EmptyUser", "", "Moldova", isDonator, 0.0},
				{"EmptyDB", "Moldova", "", isDonator, 0.0},
			}

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					user := models.UserData{UserClaim: models.UserClaim{RegCountry: c.userVal}}
					db := models.DBRecord{RegCountry: c.dbVal}

					result := countryCalc.Calculate(user, db, weights)

					assert.InDelta(t, c.expected, result.Result, 0.001, "Value mismatch in %s", c.name)

					if c.expected > 0 {
						assert.Equal(t, "match", result.Status)
					} else if c.userVal == "" || c.dbVal == "" {
						assert.Equal(t, "no_data", result.Status)
					} else {
						assert.Equal(t, "no_match", result.Status)
					}
				})
			}
		})

		// Тесты для Города
		t.Run(prefix+"_City", func(t *testing.T) {
			cases := []testCase{
				{"Match", "Chisinau", "Chisinau", isDonator, weights.RegCity},
				{"WithSpaces", " Chisinau ", "Chisinau", isDonator, weights.RegCity},
				{"Mismatch", "Chisinau", "Balti", isDonator, 0.0},
				{"EmptyBoth", "", "", isDonator, 0.0},
			}

			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					user := models.UserData{UserClaim: models.UserClaim{RegCity: c.userVal}}
					db := models.DBRecord{RegCity: c.dbVal}

					result := cityCalc.Calculate(user, db, weights)

					assert.InDelta(t, c.expected, result.Result, 0.001, "Value mismatch in %s", c.name)

					if c.expected > 0 {
						assert.Equal(t, "match", result.Status)
					} else if c.userVal == "" || c.dbVal == "" {
						assert.Equal(t, "no_data", result.Status)
					} else {
						assert.Equal(t, "no_match", result.Status)
					}
				})
			}
		})
	}
}
