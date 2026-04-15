package tests

import (
    "fmt"
    "testing"
    "time"

    "techsupport/core/internal/logic"
    "techsupport/core/pkg/models"
    "github.com/stretchr/testify/assert"
)

func TestCalculateScoreForCreationAge_RegDate_BoundaryCases(t *testing.T) {
    baseDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

    calc := logic.RegDateCalculator{Log: nil} 

    type testCase struct {
        name           string
        userDate       time.Time
        dbDate         time.Time
        isDonator      bool
        expectedResult float64
        expectedStatus string
    }

    var cases []testCase

    isDonatorOptions := []bool{false, true}
    for _, isDonator := range isDonatorOptions {
        weights := logic.GetWeights(isDonator)
        weight := weights.RegDate

        for i := -5; i <= 5; i++ {
            for j := 0; j <= 3; j++ {
                expectedRes := 0.0
                expectedStatus := "no_match"
                diffMonths := float64(j)
                
                // Твоя логика границ
                if diffMonths <= 2 {
                    expectedRes = weight 
                    expectedStatus = "match"
                } else if diffMonths <= 4 {
                    expectedRes = weight * 0.5 
                    expectedStatus = "partial"
                }

                userTime := baseDate.Add(time.Duration(i) * time.Hour).AddDate(0, j, 0)
                
                prefix := "F2P"
                if isDonator { prefix = "Solvent" }

                cases = append(cases, testCase{
                    name:           fmt.Sprintf("%s_hours=%d_months=%d", prefix, i, j),
                    userDate:       userTime,
                    dbDate:         baseDate,
                    isDonator:      isDonator,
                    expectedResult: expectedRes,
                    expectedStatus: expectedStatus,
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

            assert.InDelta(t, c.expectedResult, result.Result, 0.001, "Test failed: %s", c.name)
            
            assert.Equal(t, c.expectedStatus, result.Status, "Status mismatch: %s", c.name)
        })
    }
}