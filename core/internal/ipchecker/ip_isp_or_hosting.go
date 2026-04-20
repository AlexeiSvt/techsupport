// Package ipchecker provides tools for validating and retrieving information about IP addresses.
package ipchecker

import (
	"strings"
	logPkg "techsupport/log/pkg"
)

// IsResidential checks if the IP address belongs to a residential Internet Service Provider (ISP).
// Residential IPs are generally considered more trustworthy as they are assigned to home users.
func (info *IpApiResponse) IsResidential(log logPkg.Logger) bool {
	// Safety check: ensure the receiver is not nil to prevent panics.
	if info == nil {
		if log != nil {
			log.Warnw("IP info is nil during residential check")
		}
		return false
	}

	// Compare ASN type to "isp". Use EqualFold for case-insensitive matching.
	isRes := strings.EqualFold(info.ASN.Number, "isp")
	
	if log != nil {
		log.Debugw("IP residential check", 
			"ip", info.IP, 
			"asn_type", info.ASN.Number, 
			"is_residential", isRes,
		)
	}
	
	return isRes
}

// IsHosting checks if the IP address originates from a data center or hosting provider.
// IPs flagged as "hosting" are often associated with VPNs, proxies, or automated bots.
func (info *IpApiResponse) IsHosting(log logPkg.Logger) bool {
	// Safety check: ensure the receiver is not nil to prevent panics.
	if info == nil {
		if log != nil {
			log.Warnw("IP info is nil during hosting check")
		}
		return false
	}

	// Compare ASN type to "hosting". Use EqualFold for case-insensitive matching.
	isHosting := strings.EqualFold(info.ASN.Number, "hosting")

	if isHosting && log != nil {
		log.Infow("hosting IP detected", 
			"ip", info.IP, 
			"org", info.ASN.Org,
		)
	}

	return isHosting
}