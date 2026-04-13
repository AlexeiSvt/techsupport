package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"techsupport/core/internal/ipchecker"
	"testing"
)

func TestGetIpInfo_TrustScore(t *testing.T) {
	tests := []struct {
		name            string
		apiResponse     ipchecker.IpApiResponse
		httpStatus      int
		shouldError     bool
		expectedPenalty float64
	}{
		{
			name: "Datacenter IP (Google Public DNS)",
			apiResponse: ipchecker.IpApiResponse{
				IP:           "8.8.8.8",
				TrustScore:   60,
				IsDatacenter: true,
				ASN:          ipchecker.ASNInfo{Org: "Google LLC"},
			},
			httpStatus:      http.StatusOK,
			shouldError:     false,
			expectedPenalty: 40.0,
		},
		{
			name: "Анонимайзер (VPN)",
			apiResponse: ipchecker.IpApiResponse{
				IP:         "45.1.1.1",
				TrustScore: 5,
				IsVPN:      true,
			},
			httpStatus:      http.StatusOK,
			shouldError:     false,
			expectedPenalty: 95.0, 
		},
		{
			name: "Критический абузер (Черный список)",
			apiResponse: ipchecker.IpApiResponse{
				IP:         "1.1.1.1",
				TrustScore: 0,
			},
			httpStatus:      http.StatusOK,
			shouldError:     false,
			expectedPenalty: 100.0,
		},
		{
			name:            "Ошибка сервера API (500)",
			apiResponse:     ipchecker.IpApiResponse{},
			httpStatus:      http.StatusInternalServerError,
			shouldError:     true,
			expectedPenalty: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Api-Key") == "" && tt.httpStatus == http.StatusOK {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.WriteHeader(tt.httpStatus)
				_ = json.NewEncoder(w).Encode(tt.apiResponse)
			}))
			defer server.Close()

			// 2. Вызываем вспомогательную функцию, имитирующую запрос
			res, err := MockGetIpInfo(server.URL, tt.apiResponse.IP, "test-api-key")

			// 3. Проверка на ошибки
			if tt.shouldError {
				if err == nil {
					t.Errorf("[%s] Ожидалась ошибка, но запрос прошел", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("[%s] Неожиданная ошибка: %v", tt.name, err)
			}

			penalty := res.GetPenaltyScore()
			if penalty != tt.expectedPenalty {
				t.Errorf("[%s] Математика неверна: ждали %.2f, получили %.2f", tt.name, tt.expectedPenalty, penalty)
			}

			if tt.apiResponse.ASN.Org != "" && res.GetOperator() != tt.apiResponse.ASN.Org {
				t.Errorf("[%s] Неверный оператор: ждали %s, получили %s", tt.name, tt.apiResponse.ASN.Org, res.GetOperator())
			}
		})
	}
}

func MockGetIpInfo(baseUrl string, ip string, apiKey string) (*ipchecker.IpApiResponse, error) {
	u := fmt.Sprintf("%s?ip=%s", baseUrl, ip)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	var result ipchecker.IpApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}