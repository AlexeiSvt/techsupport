// Package tests implements integration testing for external service dependencies.
package tests

import (
	"os"
	"techsupport/core/internal/ipchecker"
	logPkg "techsupport/log/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetIpInfo_RealIntegration performs a live network check against the IP Intelligence API.
// It verifies that the enrichment layer can successfully fetch, parse, and score 
// real-world IP addresses like Google and Cloudflare DNS.
func TestGetIpInfo_RealIntegration(t *testing.T) {
	// Ensure the environment is configured for live API calls.
	apiKey := os.Getenv("API_IP_INFO_KEY")
	if apiKey == "" {
		t.Skip("SKIP: API_IP_INFO_KEY not set. Live network tests require an API key.")
	}

	// Matrix of reliable public IPs for stable integration testing.
	tests := []struct {
		name string
		ip   string
	}{
		{name: "Google DNS", ip: "8.8.8.8"},
		{name: "Cloudflare DNS", ip: "1.1.1.1"},
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // FIX: Removed context.Context because the function signature 
            // only expects (Logger, string).
            // We cast nil to the specific Logger interface to avoid ambiguity.
            var nilLog logPkg.Logger = nil 
            res, err := ipchecker.GetIpInfo(nilLog, tt.ip)

            if err != nil {
                t.Fatalf("[%s] API error: %v", tt.name, err)
            }

            assert.NotEmpty(t, res.IP)
            assert.NotEmpty(t, res.ASN.Org)

			// Calculating the penalty score based on the received IP intelligence.
			// Passing nil for logger to keep test output clean.
			penalty := res.GetPenaltyScore(nil)

			t.Logf("[%s] Final Report: IP=%s, Penalty=%.2f, Trust=%.2f", 
				tt.name, res.IP, penalty, res.TrustScore)

			// Mathematical integrity check: 
			// In our scoring model, Trust and Penalty must always be complementary to 100.
			assert.InDelta(t, 100.0, res.TrustScore+penalty, 0.01, 
				"Scoring imbalance: Trust and Penalty must sum to 100")
		})
	}
}