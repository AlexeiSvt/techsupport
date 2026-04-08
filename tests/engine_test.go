package tests

import (
    "math"
    "runtime/debug"
    "testing"
    "time"
    "os"

    "techsupport/internal/engine"
    "techsupport/internal/models"
)

func TestCalculateFinalScore_Production(t *testing.T) {
    // Чтобы ipchecker не ругался на пустой ключ и не паниковал
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
                        // Добавляем пустую транзакцию, чтобы калькулятор транзакций не паниковал
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
            // Увеличиваем таймаут для сетевых ожиданий, если они всё же есть
            defer func() {
                if r := recover(); r != nil {
                    // Если паника всё еще unimplemented, значит надо мокать GetIpInfo
                    t.Fatalf("\n🛑 КРИТИЧЕСКИЙ СБОЙ: %s\nПричина: %v\n\nСтек трейс:\n%s", tt.name, r, debug.Stack())
                }
            }()

            // ВАЖНО: Если ipchecker всё равно валит тест, 
            // в реальной разработке GetIpInfo выносится в интерфейс.
            got := engine.CalculateFinalScore(tt.input)

            if math.Abs(got-tt.expected) > 0.01 {
                t.Errorf("\n❌ ОШИБКА РАСЧЕТА [%s]:\nОжидали: %.2f\nПолучили: %.2f", tt.name, tt.expected, got)
            }
        })
    }
}