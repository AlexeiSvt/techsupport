// Package ipchecker provides tools for validating and retrieving information about IP addresses.
package ipchecker

import (
	"net/http"
	"time"
)

// ASNInfo represents Autonomous System Number information.
// It identifies the organization owning the IP range and the type of network.
type ASNInfo struct {
	Org  string `json:"org"`  // Example: "Google LLC", "Comcast Cable"
	Type string `json:"type"` // Example: "isp", "hosting", "business", "edu"
}

// IpApiResponse is the data transfer object (DTO) for the external IP verification API.
// It contains reputation metrics and infrastructure flags.
type IpApiResponse struct {
	IP           string  `json:"ip"`
	TrustScore   float64 `json:"trust_score"`   // Reputation score from 0 to 100
	IsMobile     bool    `json:"is_mobile"`     // True if the IP belongs to a cellular carrier
	IsDatacenter bool    `json:"is_datacenter"` // True if the IP is from a known data center
	IsVPN        bool    `json:"is_vpn"`        // True if the IP is identified as a VPN/Proxy
	ASN          ASNInfo `json:"asn"`
}

var (
	// httpClient is a globally shared, thread-safe HTTP client.
	// It is pre-configured with optimized connection pooling for high-concurrency environments.
	httpClient = &http.Client{
		// Timeout includes connection, any redirects, and reading the response body.
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			// MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts.
			MaxIdleConns: 100,
			// MaxIdleConnsPerHost prevents a single host from exhausting the entire connection pool.
			MaxIdleConnsPerHost: 20,
			// IdleConnTimeout defines how long an idle connection stays open before closing.
			IdleConnTimeout: 90 * time.Second,
		},
	}
)