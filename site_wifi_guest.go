package unifi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type WifiGuestConfig struct {
	UploadSpeed    int    // upload speed in Kbps
	DownloadSpeed  int    // download speed in Kbps
	TransferLimit  int    // data transfer limits in MB
	AccessPointMac string // Access Point MAC address to which the client is connected, should result in faster authorization
}

// AuthorizeWiFiGuest will authorize a wifi guest on the network
// site - site to allow the guest
// mac - client mac to authorize
// duration - time for wifi authorization, if <=0 , then it will default to 1hr
// wifiGuestConfig - optional parameters to limit the client
func (c *Client) AuthorizeWiFiGuest(site string, mac string, duration time.Duration, wifiGuestConfig *WifiGuestConfig) (*GenericResponse, error) {
	if duration.Minutes() <= 0 {
		duration = time.Hour * 1
	}

	payload := map[string]interface{}{
		"cmd":     "authorize-guest",
		"mac":     strings.ToLower(mac),
		"minutes": int(duration.Minutes()),
	}
	// optionally each parameter if non-zero
	if wifiGuestConfig != nil {
		if wifiGuestConfig.UploadSpeed >= 0 {
			payload["up"] = wifiGuestConfig.UploadSpeed
		}

		if wifiGuestConfig.DownloadSpeed >= 0 {
			payload["down"] = wifiGuestConfig.DownloadSpeed
		}

		if wifiGuestConfig.TransferLimit >= 0 {
			payload["bytes"] = wifiGuestConfig.TransferLimit
		}

		if wifiGuestConfig.AccessPointMac != "" {
			payload["ap_mac"] = wifiGuestConfig.AccessPointMac
		}
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/stamgr", bytes.NewReader(data), &resp)
	return &resp, err
}

func (c *Client) UnAuthorizeWiFiGuest(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "unauthorize-guest",
		"mac": strings.ToLower(mac),
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/stamgr", bytes.NewReader(data), &resp)
	return &resp, err
}
