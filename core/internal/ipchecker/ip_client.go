// Package ipchecker provides tools for validating and retrieving information about IP addresses.
// It integrates with external APIs to fetch geolocation and reputation data.
package ipchecker

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"techsupport/core/internal/constants"
	iperror "techsupport/core/internal/ipchecker/ip_errors"
	logPkg "techsupport/log/pkg"
)

// apiBaseURL is the endpoint for the external IP information service.
var apiBaseURL = constants.ApiBaseURL

// GetIpInfoWithContext fetches detailed information about a specific IP address.
// It respects the provided context for timeouts and cancellation, ensuring 
// the application doesn't hang on slow network responses.
func GetIpInfoWithContext(ctx context.Context, log logPkg.Logger, ip string) (*IpApiResponse, error) {
	// Clean input to avoid formatting issues.
	ip = strings.TrimSpace(ip)

	// Basic validation: Check if IP is provided.
	if ip == "" {
		return nil, iperror.ErrIPEmpty
	}
	
	// Format validation: Ensure the string is a valid IPv4 or IPv6 address.
	if net.ParseIP(ip) == nil {
		if log != nil {
			log.Warnw("invalid ip format", "ip", ip)
		}
		return nil, fmt.Errorf("%w: %s", iperror.ErrIPInvalidFormat, ip)
	}

	// Security: Retrieve the API key from a secure source/environment.
	apiKey := getAPIKey()
	if apiKey == "" {
		if log != nil {
			log.Errorw("critical config missing", "env", "API_IP_INFO_KEY")
		}
		return nil, iperror.ErrApiKeyMissing
	}

	// URL Construction: Safely build the request URL with query parameters.
	u, err := url.Parse(apiBaseURL)
	if err != nil {
		if log != nil {
			log.Errorw("failed to parse API base URL", "url", apiBaseURL, "err", err)
		}
		return nil, fmt.Errorf("%w: %v", iperror.ErrApiUrlInvalid, err)
	}

	q := u.Query()
	q.Set("ip", ip)
	u.RawQuery = q.Encode()

	// Request creation: Linking the request to the lifecycle of the Context.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	// Authorization: Set the required API key in the headers.
	req.Header.Set("X-Api-Key", apiKey)
	
	if log != nil {
		log.Debugw("requesting IP info", "url", u.String())
	}

	// Network Execution: Perform the actual HTTP call via the pre-configured client.
	resp, err := httpClient.Do(req)
	if err != nil {
		if log != nil {
			log.Errorw("external IP API request failed", "ip", ip, "err", err)
		}
		return nil, fmt.Errorf("%w: %v", iperror.ErrApiRequest, err)
	}
	// Important: Always close the body to prevent memory leaks and keep connections reusable.
	defer resp.Body.Close()

	// Protocol validation: Check if the API returned a successful HTTP status.
	if resp.StatusCode != http.StatusOK {
		if log != nil {
			log.Errorw("external IP API returned error status", "status", resp.StatusCode, "ip", ip)
		}
		return nil, fmt.Errorf("%w: %d", iperror.ErrApiStatus, resp.StatusCode)
	}

	// Data Decoding: Convert the JSON response body into a Go struct.
	var result IpApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		if log != nil {
			log.Errorw("failed to decode IP API response", "err", err)
		}
		return nil, fmt.Errorf("%w: %v", iperror.ErrApiDecode, err)
	}

	// Sanity Check: Ensure the API actually returned data for the requested IP.
	if result.IP == "" {
		if log != nil {
			log.Warnw("IP API returned empty result object", "ip", ip)
		}
		return nil, iperror.ErrApiResponseEmpty
	}

	if log != nil {
		log.Infow("IP info retrieved successfully", "ip", ip)
	}

	return &result, nil
}