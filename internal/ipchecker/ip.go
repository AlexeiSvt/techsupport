package ipchecker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var API_IP_INFO_KEY = os.Getenv("API_IP_INFO_KEY")

type IpApiResponse struct {
	IP           string `json:"ip"`
	IsBogon      bool   `json:"is_bogon"`
	IsMobile     bool   `json:"is_mobile"`
	IsSatellite  bool   `json:"is_satellite"`
	IsCrawler    bool   `json:"is_crawler"`
	IsDatacenter bool   `json:"is_datacenter"`
	IsTor        bool   `json:"is_tor"`
	IsProxy      bool   `json:"is_proxy"`
	IsVPN        bool   `json:"is_vpn"`
	IsAbuser     bool   `json:"is_abuser"`
	ASN          struct {
		Org  string `json:"org"`
		Type string `json:"type"`
	} `json:"asn"`
}

func (info *IpApiResponse) GetPenaltyScore() float64 {

    isDirty := info.IsDatacenter || info.IsVPN || info.IsProxy || info.IsTor || info.IsAbuser || info.IsBogon

    if !isDirty || info.IsMobile || info.IsSatellite {
        return 0
    }

    return 100
}

func GetIpInfo(ip string) (*IpApiResponse, error) {
	url := fmt.Sprintf("https://api.ipapi.is?ip=%s&key=%s", ip, API_IP_INFO_KEY)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", API_IP_INFO_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() 

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var result IpApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}