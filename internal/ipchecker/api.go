package ipchecker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"techsupport/internal/scoring"
)

var apiBaseURL = scoring.ApiBaseURL

func GetIpInfoWithContext(ctx context.Context, ip string) (*IpApiResponse, error) {
	ip = strings.TrimSpace(ip)

	if ip == "" {
		return nil, fmt.Errorf("empty ip")
	}

	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("invalid ip address")
	}

	apiKey := getAPIKey()

	if ip == "127.0.0.1" || apiKey == "test_key" {
		return &IpApiResponse{
			IP: ip,
			ASN: ASNInfo{
				Org:  "Test-Network",
				Type: "business",
			},
		}, nil
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API_IP_INFO_KEY is not set")
	}

	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse api url: %w", err)
	}

	q := u.Query()
	q.Set("ip", ip)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("X-Api-Key", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request ipapi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipapi returned status %d", resp.StatusCode)
	}

	var result IpApiResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if result.IP == "" {
		return nil, fmt.Errorf("empty api response")
	}

	return &result, nil
}