package ipchecker

import (
	"strings"
	logPkg "techsupport/log/pkg"
)

func (info *IpApiResponse) IsResidential(log logPkg.Logger) bool {
	if info == nil {
		log.Warnw("IP info is nil during residential check")
		return false
	}

	isRes := strings.EqualFold(info.ASN.Type, "isp")
	
	log.Debugw("IP residential check", 
		"ip", info.IP, 
		"asn_type", info.ASN.Type, 
		"is_residential", isRes,
	)
	
	return isRes
}

func (info *IpApiResponse) IsHosting(log logPkg.Logger) bool {
	if info == nil {
		log.Warnw("IP info is nil during hosting check")
		return false
	}

	isHosting := strings.EqualFold(info.ASN.Type, "hosting")

	if isHosting {
		log.Infow("hosting IP detected", "ip", info.IP, "org", info.ASN.Org)
	}

	return isHosting
}