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
        return float64(constants.FullPenalty)
    }

    penalty := 100.0 - info.TrustScore

    if penalty > float64(constants.FullPenalty) {
        penalty = float64(constants.FullPenalty)
    }

    if penalty < 0 {
        penalty = 0
    }

    return penalty
}

func (info *IpApiResponse) GetOperator() string {
    if info == nil {
        return "Unknown"
    }
    return info.ASN.Org
}

func GetIpInfo(ip string) (*IpApiResponse, error) {
    return GetIpInfoWithContext(context.Background(), ip)
}