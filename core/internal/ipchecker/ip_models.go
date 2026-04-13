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
    TrustScore   float64 `json:"trust_score"`
    IsMobile     bool    `json:"is_mobile"`
    IsDatacenter bool    `json:"is_datacenter"`
    IsVPN        bool    `json:"is_vpn"`
    ASN          ASNInfo `json:"asn"`
}

var (
    httpClient = &http.Client{
        Timeout: 5 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 20,
            IdleConnTimeout:     90 * time.Second,
        },
    }
)