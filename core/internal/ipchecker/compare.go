package ipchecker

import "strings"

func (info *IpApiResponse) IsResidential() bool {
	return info != nil && strings.EqualFold(info.ASN.Type, "isp")
}

func (info *IpApiResponse) IsHosting() bool {
	return info != nil && strings.EqualFold(info.ASN.Type, "hosting")
}