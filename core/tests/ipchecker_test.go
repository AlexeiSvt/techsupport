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
			name: "Clear home user",
			apiResponse: ipchecker.IpApiResponse{
				IP: "1.1.1.1",
				IsDatacenter: false,
				IsMobile:     false,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "Mobile user (Trusted IP)",
			apiResponse: ipchecker.IpApiResponse{
				IP: "8.8.8.8",
				IsMobile: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "VPN detected",
			apiResponse: ipchecker.IpApiResponse{
				IP: "45.1.1.1",
				IsVPN: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "Critical Abuser",
			apiResponse: ipchecker.IpApiResponse{
				IP: "66.6.6.6",
				IsAbuser: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
		{
			name: "API Error (401 Unauthorized)",
			apiResponse: ipchecker.IpApiResponse{},
			httpStatus:  http.StatusUnauthorized,
			shouldError: true,
		},
		{
			name: "Bogon (non-existent IP)",
			apiResponse: ipchecker.IpApiResponse{
				IsBogon: true,
			},
			httpStatus:  http.StatusOK,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.httpStatus)
				json.NewEncoder(w).Encode(tt.apiResponse)
			}))
			defer server.Close()

			
			
			res, err := GetIpInfoFromUrl(server.URL, "1.1.1.1")

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error, but none occurred")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if res.IsDatacenter != tt.apiResponse.IsDatacenter {
				t.Errorf("Expected IsDatacenter=%v, got %v", tt.apiResponse.IsDatacenter, res.IsDatacenter)
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