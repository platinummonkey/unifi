package unifi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SiteRougeAccessPoint struct {
	ID              string      `json:"_id"`
	Age             int         `json:"age"`
	AccessPointMAC  string      `json:"ap_mac"`
	Band            string      `json:"band"`
	BSSID           string      `json:"bssid"`
	BW              int         `json:"bw"`
	CenterFrequency int         `json:"center_freq"`
	Channel         int         `json:"channel"`
	ESSID           string      `json:"essid"`
	Frequency       int         `json:"freq"`
	IsAdHoc         bool        `json:"is_adhoc"`
	IsRogue         bool        `json:"is_rogue"`
	IsUbnt          bool        `json:"is_ubnt"`
	LastSeen        interface{} `json:"last_seen"`
	Noise           int         `json:"noise"`
	OUI             string      `json:"oui"`
	Radio           string      `json:"radio"`
	RadioName       string      `json:"radio_name"`
	ReportTime      int64       `json:"report_time"`
	RSSI            int         `json:"rssi"`
	RSSIAge         int         `json:"rssi_age"`
	Security        string      `json:"security"`
	Signal          int         `json:"signal"`
	SiteID          string      `json:"site_id"`
}

type SiteRougeAccessPointResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data []SiteRougeAccessPoint `json:"data"`
}

func (c *Client) SiteRougeAccessPoints(site string, seenWithinHours int) (*SiteRougeAccessPointResponse, error) {
	if seenWithinHours < 0 {
		seenWithinHours = 24
	}

	payload := map[string]interface{}{
		"within": seenWithinHours,
	}
	data, _ := json.Marshal(payload)

	var resp SiteRougeAccessPointResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/rogueap", bytes.NewReader(data), &resp)
	return &resp, err
}
