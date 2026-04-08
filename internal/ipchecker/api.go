package ipchecker

import (
	"net/http"
	"time"
)

type ASNInfo struct {
	Org  string `json:"org"`
	Type string `json:"type"`
}

type IpApiResponse struct {
	IP           string  `json:"ip"`
	IsBogon      bool    `json:"is_bogon"`
	IsMobile     bool    `json:"is_mobile"`
	IsSatellite  bool    `json:"is_satellite"`
	IsCrawler    bool    `json:"is_crawler"`
	IsDatacenter bool    `json:"is_datacenter"`
	IsTor        bool    `json:"is_tor"`
	IsProxy      bool    `json:"is_proxy"`
	IsVPN        bool    `json:"is_vpn"`
	IsAbuser     bool    `json:"is_abuser"`
	ASN          ASNInfo `json:"asn"`
}

var (
	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
) 
