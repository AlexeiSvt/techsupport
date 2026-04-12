package tests

import (
    "fmt"
    "testing"
    "time"

    "techsupport/core/internal/logic"
    "techsupport/core/internal/models"
    "github.com/stretchr/testify/assert"
)

func TestCalculateScoreForCreationAge_RegDate_BoundaryCases(t *testing.T) {
    baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    
    calc := logic.RegDateCalculator{}

    type testCase struct {
        name      string
        userDate  time.Time
        dbDate    time.Time
        isDonator bool
        expected  float64
    }

    var cases []testCase

    isDonatorOptions := []bool{false, true}
    for _, isDonator := range isDonatorOptions {
        weight := 12.5
        if isDonator {
            weight = 7.5
        }

        for i := -5; i <= 5; i++ {
            for j := 0; j <= 3; j++ {
                expected := 0.0
                diffMonths := float64(j)
                
                if diffMonths <= 2 {
                    expected = weight 
                } else if diffMonths <= 4 {
                    expected = weight * 0.5 
                }

                userTime := baseDate.Add(time.Duration(i) * time.Hour).AddDate(0, j, 0)
                
                prefix := "F2P"
                if isDonator { prefix = "Solvent" }

                cases = append(cases, testCase{
                    name:      fmt.Sprintf("%s_hours=%d_months=%d", prefix, i, j),
                    userDate:  userTime,
                    dbDate:    baseDate,
                    isDonator: isDonator,
                    expected:  expected,
                })
            }
        }
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            weights := logic.GetWeights(c.isDonator)
            
            user := models.UserData{
                UserClaim: models.UserClaim{RegDate: c.userDate},
            }
            db := models.DBRecord{
                RegDate: c.dbDate,
            }

            result := calc.Calculate(user, db, weights)

            assert.InDelta(t, c.expected, result.Result, 0.001, "Test failed: %s", c.name)
                        if c.expected > 0 && c.expected == weights.RegDate {
                assert.Equal(t, "match", result.Status)
            }
        })
    }
}