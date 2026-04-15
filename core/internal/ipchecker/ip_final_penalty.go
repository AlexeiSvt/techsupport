// Package ipchecker provides tools for validating and retrieving information about IP addresses.
package ipchecker

import (
	"context"
	"os"
	"strings"
	"techsupport/core/internal/constants"
	logPkg "techsupport/log/pkg"
)

// getAPIKey retrieves the external service authentication token from environment variables.
// It trims any accidental whitespace to prevent authentication failures.
func getAPIKey() string {
	return strings.TrimSpace(os.Getenv("API_IP_INFO_KEY"))
}

// GetPenaltyScore converts the IP's TrustScore into a penalty value for the scoring engine.
// If the IP info is missing (nil), it returns a maximum penalty as a fail-safe measure.
// The penalty logic follows the formula: Penalty = 100 - TrustScore.
func (info *IpApiResponse) GetPenaltyScore(log logPkg.Logger) float64 {
	// Security/Stability: Handle cases where IP data could not be fetched.
	if info == nil {
		if log != nil {
			log.Warnw("IP info is nil, applying full penalty due to lack of data")
		}
		return float64(constants.FullPenalty)
	}

	// Logic: Lower trust results in a higher penalty.
	penalty := 100.0 - info.TrustScore

	// Constraint: Ensure the penalty does not exceed the globally defined maximum.
	if penalty > float64(constants.FullPenalty) {
		penalty = float64(constants.FullPenalty)
	}

	// Constraint: Ensure penalty is never negative.
	if penalty < 0 {
		penalty = 0
	}

	if log != nil {
		log.Debugw("calculated IP penalty score", 
			"ip", info.IP, 
			"trust_score", info.TrustScore, 
			"final_penalty", penalty,
		)
	}

	return penalty
}

// GetOperator returns the name of the Internet Service Provider (ISP) or Organization 
// associated with the IP's Autonomous System Number (ASN).
func (info *IpApiResponse) GetOperator() string {
	if info == nil {
		return "Unknown"
	}
	// Org usually contains the name of the ISP/Hosting provider (e.g., "Google LLC", "Comcast").
	return info.ASN.Org
}

// GetIpInfo is a convenience wrapper for GetIpInfoWithContext using a background context.
// Use this only for legacy calls or where strict timeout management is handled elsewhere.
func GetIpInfo(log logPkg.Logger, ip string) (*IpApiResponse, error) {
	return GetIpInfoWithContext(context.Background(), log, ip)
}