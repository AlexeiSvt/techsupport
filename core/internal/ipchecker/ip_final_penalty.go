package ipchecker

import (
	"context"
	"os"
	"strings"
	"techsupport/core/internal/constants"
)

func getAPIKey() string {
	return strings.TrimSpace(os.Getenv("API_IP_INFO_KEY"))
}

func (info *IpApiResponse) GetPenaltyScore() float64 {
	if info == nil {
		return constants.FullPenalty
	}

	score := 0.0

	switch {
	case info.IsBogon:
		return constants.FullPenalty
	case info.IsAbuser:
		return constants.FullPenalty
	case info.IsTor:
		return constants.FullPenalty
	case info.IsCrawler:
		return constants.FullPenalty
	}

	if info.IsDatacenter {
		score += constants.ForDatacenter
	}

	if info.IsVPN {
		score += constants.ForVPN
	}

	if info.IsProxy {
		score += constants.ForProxy
	}

	if strings.EqualFold(info.ASN.Type, "hosting") {
		score += constants.ForHosting
	}

	if strings.EqualFold(info.ASN.Type, "mobile") {
		score *= constants.PartialMatch
	}

	if info.IsMobile {
		score *= constants.PartialMatch
	}

	if info.IsSatellite {
		score *= constants.MostlyMatch
	}

	if score > constants.FullPenalty {
		score = constants.FullPenalty
	}

	return score
}


func GetIpInfo(ip string) (*IpApiResponse, error) {
	return GetIpInfoWithContext(context.Background(), ip)
}