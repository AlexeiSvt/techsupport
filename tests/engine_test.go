package tests

import (
	"math"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"techsupport/internal/engine"
	"techsupport/internal/models"
)

func TestCalculateFinalScore_Production(t *testing.T) {
	// Устанавливаем дефолтный ключ, если его нет в системе
	if os.Getenv("API_IP_INFO_KEY") == "" {
		os.Setenv("API_IP_INFO_KEY", "test_vortex_key")
	}

	now := time.Now()

	// Таблица тестов
	tests := []struct {
		name     string
		input    models.InputData
		expected float64
	}{
		{
			name: "01. Full Match (Clean IP)",
			input: models.InputData{
				UserData: models.UserData{
					IP: "127.0.0.1",
					UserClaim: models.UserClaim{
						AccTag: "Vortex#1", RegCountry: "RU", RegCity: "MSK",
						FirstEmail: "v@v.com", Phone: "799", FirstDevice: "PC",
						Devices: []string{"PC"}, RegDate: now,
					},
				},
				DBRecord: models.DBRecord{
					AccTag: "Vortex#1", RegCountry: "RU", RegCity: "MSK",
					FirstEmail: "v@v.com", Phone: "799", FirstDevice: "PC",
					Devices: []string{"PC"}, RegDate: now,
				},
			},
			expected: 100.0,
		},
		// Здесь добавь остальные 29 кейсов...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Отлов паники (unimplemented и прочее)
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("\n🔥 ПАНИКА: %s\nПричина: %v\nСтек:\n%s", tt.name, r, debug.Stack())
				}
			}()

			got := engine.CalculateFinalScore(tt.input)

			// Сравнение float с допуском
			if math.Abs(got-tt.expected) > 0.01 {
				t.Errorf("\n❌ %s: Ожидали %.2f, получили %.2f", tt.name, tt.expected, got)
			}
		})
	}
}