package ipchecker

import (
    "context"
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "net/url"
    "strings"
    "techsupport/core/internal/constants"
    iperror "techsupport/core/internal/ipchecker/ip_errors"
    logPkg "techsupport/log/pkg"
)

var apiBaseURL = constants.ApiBaseURL

func GetIpInfoWithContext(ctx context.Context, log logPkg.Logger, ip string) (*IpApiResponse, error) {
    ip = strings.TrimSpace(ip)

    if ip == "" {
        return nil, iperror.ErrIPEmpty
    }
    
    if net.ParseIP(ip) == nil {
        if log != nil {
            log.Warnw("invalid ip format", "ip", ip)
        }
        return nil, fmt.Errorf("%w: %s", iperror.ErrIPInvalidFormat, ip)
    }

    apiKey := getAPIKey()
    if apiKey == "" {
        if log != nil {
            log.Errorw("critical config missing", "env", "API_IP_INFO_KEY")
        }
        return nil, iperror.ErrApiKeyMissing
    }

    u, err := url.Parse(apiBaseURL)
    if err != nil {
        if log != nil {
            log.Errorw("failed to parse API base URL", "url", apiBaseURL, "err", err)
        }
        return nil, fmt.Errorf("%w: %v", iperror.ErrApiUrlInvalid, err)
    }

    q := u.Query()
    q.Set("ip", ip)
    u.RawQuery = q.Encode()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("build request: %w", err)
    }

    req.Header.Set("X-Api-Key", apiKey)
    
    if log != nil {
        log.Debugw("requesting IP info", "url", u.String())
    }

    resp, err := httpClient.Do(req)
    if err != nil {
        if log != nil {
            log.Errorw("external IP API request failed", "ip", ip, "err", err)
        }
        return nil, fmt.Errorf("%w: %v", iperror.ErrApiRequest, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        if log != nil {
            log.Errorw("external IP API returned error status", "status", resp.StatusCode, "ip", ip)
        }
        return nil, fmt.Errorf("%w: %d", iperror.ErrApiStatus, resp.StatusCode)
    }

    var result IpApiResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        if log != nil {
            log.Errorw("failed to decode IP API response", "err", err)
        }
        return nil, fmt.Errorf("%w: %v", iperror.ErrApiDecode, err)
    }

    if result.IP == "" {
        if log != nil {
            log.Warnw("IP API returned empty result object", "ip", ip)
        }
        return nil, iperror.ErrApiResponseEmpty
    }

    if log != nil {
        log.Infow("IP info retrieved successfully", "ip", ip)
    }

    return &result, nil
}