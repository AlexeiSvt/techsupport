package tests

import (
	"os"
	"techsupport/core/internal/ipchecker"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIpInfo_RealIntegration(t *testing.T) {
	apiKey := os.Getenv("API_IP_INFO_KEY")
	if apiKey == "" {
		t.Skip("SKIP: API_IP_INFO_KEY not set")
	}

	tests := []struct {
		name string
		ip   string
	}{
		{name: "Google DNS", ip: "8.8.8.8"},
		{name: "Cloudflare DNS", ip: "1.1.1.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ipchecker.GetIpInfo(nil, tt.ip)

			if err != nil {
				t.Fatalf("[%s] API error: %v", tt.name, err)
			}

			assert.NotEmpty(t, res.IP)
			assert.NotEmpty(t, res.ASN.Org)

			penalty := res.GetPenaltyScore(nil)

			t.Logf("[%s] Result: IP=%s, Penalty=%.2f, Trust=%.2f", 
				tt.name, res.IP, penalty, res.TrustScore)

			assert.InDelta(t, 100.0, res.TrustScore+penalty, 0.01)
		})
	}
}