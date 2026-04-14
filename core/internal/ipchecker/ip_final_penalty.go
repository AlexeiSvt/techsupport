package ipchecker

import (
    "context"
    "os"
    "strings"
    "techsupport/core/internal/constants"
    logPkg "techsupport/log/pkg"
)

func getAPIKey() string {
    return strings.TrimSpace(os.Getenv("API_IP_INFO_KEY"))
}

func (info *IpApiResponse) GetPenaltyScore(log logPkg.Logger) float64 {
    if info == nil {
        if log != nil {
            log.Warnw("IP info is nil, applying full penalty")
        }
        return float64(constants.FullPenalty)
    }

    penalty := 100.0 - info.TrustScore

    if penalty > float64(constants.FullPenalty) {
        penalty = float64(constants.FullPenalty)
    }

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

func (info *IpApiResponse) GetOperator() string {
    if info == nil {
        return "Unknown"
    }
    return info.ASN.Org
}

func GetIpInfo(log logPkg.Logger, ip string) (*IpApiResponse, error) {
    return GetIpInfoWithContext(context.Background(), log, ip)
}