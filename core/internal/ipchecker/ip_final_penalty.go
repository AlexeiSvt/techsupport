package ipchecker

import (
	"context"
	"os"
	"strings"
	"techsupport/core/internal/scoring"
)

func getAPIKey() string {
	return strings.TrimSpace(os.Getenv("API_IP_INFO_KEY"))
}

func (info *IpApiResponse) GetPenaltyScore() float64 {
	if info == nil {
		// Fail closed: external IP reputation unavailable, treat as maximum risk
		return scoring.FullPenalty
	}

	score := 0.0

	switch {
	case info.IsBogon:
		return scoring.FullPenalty
	case info.IsAbuser:
		return scoring.FullPenalty
	case info.IsTor:
		return scoring.FullPenalty
	case info.IsCrawler:
		return scoring.FullPenalty
	}

	if info.IsDatacenter {
		score += scoring.ForDatacenter
	}

	if info.IsVPN {
		score += scoring.ForVPN
	}

	if info.IsProxy {
		score += scoring.ForProxy
	}

	if strings.EqualFold(info.ASN.Type, "hosting") {
		score += scoring.ForHosting
	}

	if strings.EqualFold(info.ASN.Type, "mobile") {
		score *= scoring.PartialMatch
	}

	if info.IsMobile {
		score *= scoring.PartialMatch
	}

	if info.IsSatellite {
		score *= scoring.MostlyMatch
	}

	if score > scoring.FullPenalty {
		score = scoring.FullPenalty
	}

	return score
}


func GetIpInfo(ip string) (*IpApiResponse, error) {
	return GetIpInfoWithContext(context.Background(), ip)
}