package tests

import (
    "math"
    "os"
    "runtime/debug"
    "strconv"
    "strings"
    "testing"
    "time"

    "techsupport/core/internal/engine"
    "techsupport/core/pkg/models"
)

func TestCalculateFinalScore_Production(t *testing.T) {
    apiKey := os.Getenv("API_IP_INFO_KEY")
    if apiKey == "" {
        t.Skip("SKIP: API_IP_INFO_KEY is not set")
    }

    now := time.Now().Add(-24 * time.Hour)

    tests := []struct {
        name     string
        input    models.InputData
        expected float64
    }{
        {
            name: "01. Full Match (Localhost - Bogon IP)",
            input: models.InputData{
                UserData: models.UserData{
                    IP: "127.0.0.1",
                    UserClaim: models.UserClaim{
                        AccTag:      "ALEXEY_DEV",
                        RegCountry:  "MD",
                        RegCity:     "Chisinau",
                        FirstEmail:  "test@mail.com",
                        Phone:       "37360000000",
                        FirstDevice: "PC",
                        Devices:     []string{"PC"},
                        RegDate:     now,
                    },
                },
                DBRecord: models.DBRecord{
                    IsDonator:   false,
                    RegCountry:  "MD",
                    RegCity:     "Chisinau",
                    FirstEmail:  "test@mail.com",
                    Phone:       "37360000000",
                    FirstDevice: "PC",
                    Devices:     []string{"PC"},
                    RegDate:     now,
                },
            },
            expected: 0.00, 
        },
        {
            name: "02. Google DNS (Datacenter)",
            input: models.InputData{
                UserData: models.UserData{
                    IP: "8.8.8.8",
                    UserClaim: models.UserClaim{
                        AccTag:      "ALEXEY_DEV",
                        RegCountry:  "US",
                        RegCity:     "Mountain View",
                        FirstEmail:  "test@mail.com",
                        Phone:       "37360000000",
                        FirstDevice: "PC",
                        Devices:     []string{"PC"},
                        RegDate:     now,
                    },
                },
                DBRecord: models.DBRecord{
                    IsDonator:   false,
                    RegCountry:  "US",
                    RegCity:     "Mountain View",
                    FirstEmail:  "test@mail.com",
                    Phone:       "37360000000",
                    FirstDevice: "PC",
                    Devices:     []string{"PC"},
                    RegDate:     now,
                },
            },
            expected: 0.00, 
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil {
                    t.Fatalf("\nCritical Failure in [%s]: %v\nStack Trace:\n%s", tt.name, r, debug.Stack())
                }
            }()

            got := engine.CalculateFinalScore(nil, tt.input)

            cleanPerc := strings.TrimSuffix(got.FinaPercentage, "%")
            gotValue, err := strconv.ParseFloat(cleanPerc, 64)
            if err != nil {
                t.Fatalf("Failed to parse final percentage: %s", got.FinaPercentage)
            }

            if math.Abs(gotValue-tt.expected) > 0.01 {
                t.Errorf("\nCalculation Error [%s]:\nExpected: %.2f\nGot: %s\nDetails: %+v",
                    tt.name, tt.expected, got.FinaPercentage, got.Details)
            }

            if len(got.Details) == 0 {
                t.Errorf("Error [%s]: Details slice is empty", tt.name)
            }
        })
    }
}