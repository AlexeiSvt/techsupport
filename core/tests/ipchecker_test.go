package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"techsupport/core/internal/ipchecker"
	"testing"
)

func TestGetIpInfo(t *testing.T) {
	tests := []struct {
		name           string
		apiResponse    ipchecker.IpApiResponse
		httpStatus     int
		shouldError    bool
		expectedRisk   string
	}{
		{
			name: "Чистая домашняя сеть",
			apiResponse: ipchecker.IpApiResponse{
				IP: "1.1.1.1",
				IsDatacenter: false,
				IsMobile:     false,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "Мобильный юзер (Траст)",
			apiResponse: ipchecker.IpApiResponse{
				IP: "8.8.8.8",
				IsMobile: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "VPN детектед",
			apiResponse: ipchecker.IpApiResponse{
				IP: "45.1.1.1",
				IsVPN: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "Критичный Абьюзер",
			apiResponse: ipchecker.IpApiResponse{
				IP: "66.6.6.6",
				IsAbuser: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "Ошибка API (401 Unauthorized)",
			apiResponse: ipchecker.IpApiResponse{},
			httpStatus:  http.StatusUnauthorized,
			shouldError: true,
		},
		{
			name: "Богон (несуществующий IP)",
			apiResponse: ipchecker.IpApiResponse{
				IsBogon: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Создаем фейковый сервер
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.httpStatus)
				json.NewEncoder(w).Encode(tt.apiResponse)
			}))
			defer server.Close()

			
			
			res, err := GetIpInfoFromUrl(server.URL, "1.1.1.1")

			if tt.shouldError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её нет")
				}
				return
			}

			if err != nil {
				t.Errorf("Неожиданная ошибка: %v", err)
			}

			if res.IsDatacenter != tt.apiResponse.IsDatacenter {
				t.Errorf("Ожидалось IsDatacenter=%v, получили %v", tt.apiResponse.IsDatacenter, res.IsDatacenter)
			}
		})
	}
}


func GetIpInfoFromUrl(baseUrl string, ip string) (*ipchecker.IpApiResponse, error) {
	url := fmt.Sprintf("%s?ip=%s", baseUrl, ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var result ipchecker.IpApiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}