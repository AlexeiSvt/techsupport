package tests

import (
    "math"
    "runtime/debug"
    "testing"
    "time"
    "os"
    "techsupport/core/internal/engine"
    "techsupport/core/internal/models"
)

func TestCalculateFinalScore_Production(t *testing.T) {
    os.Setenv("API_IP_INFO_KEY", "test_key")
    defer os.Unsetenv("API_IP_INFO_KEY")

    now := time.Now()

    tests := []struct {
        name     string
        input    models.InputData
        expected float64
    }{
        {
            name: "01. Full Match (Local Trust IP)",
            input: models.InputData{
                UserData: models.UserData{
                    IP: "127.0.0.1",
                    UserClaim: models.UserClaim{
                        RegCountry:  "RU",
                        RegCity:     "MSK",
                        FirstEmail:  "test@mail.com",
                        Phone:       "79991234567",
                        FirstDevice: "PC",
                        Devices:     []string{"PC"},
                        RegDate:     now,
                        FirstTransaction: models.Transaction{
                            Amount: 0,
                        },
                    },
                },
                DBRecord: models.DBRecord{
                    RegCountry:  "RU",
                    RegCity:     "MSK",
                    FirstEmail:  "test@mail.com",
                    Phone:       "79991234567",
                    FirstDevice: "PC",
                    Devices:     []string{"PC"},
                    RegDate:     now,
                    IsDonator:   false,
                },
            },
            expected: 100.0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil {
                    t.Fatalf("\nCritical Error: %s\n Caused by: %v\n\nStack Trace:\n%s", tt.name, r, debug.Stack())
                }
            }()

            got := engine.CalculateFinalScore(tt.input)

            if math.Abs(got-tt.expected) > 0.01 {
                t.Errorf("\nCalculation Error [%s]:\nExpected: %.2f\nGot: %.2f", tt.name, tt.expected, got)
            }
        })
    }
}