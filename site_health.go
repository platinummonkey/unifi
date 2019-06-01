package unifi

import (
	"net/http"
)

// SiteHealthData contains the site health data
type SiteHealthData struct {
	NumberAdopted      int    `json:"num_adopted"`
	NumberAccessPoints int    `json:"num_ap"`
	NumberDisabled     int    `json:"num_disabled"`
	NumberDisconnected int    `json:"num_disconnected"`
	NumberGuest        int    `json:"num_guest"`
	NumberIOTDevices   int    `json:"num_iot"`
	NumberSW           int    `json:"num_sw"`
	NumberPending      int    `json:"num_pending"`
	NumberUser         int    `json:"user"`
	RXBytesR           int64  `json:"rx_bytes-r"`
	Status             string `json:"status"`
	SubSystem          string `json:"subsystem"`
	TXBytesR           int64  `json:"tx_bytes-r"`
	Uptime             int64  `json:"uptime"`

	Gateways           []string                       `json:"gateways"`
	GatewayMAC         string                         `json:"gw_mac"`
	GatewayName        string                         `json:"gw_name"`
	GatewaySystemStats SitesVerboseGatewaySystemStats `json:"gw_system-stats"`
	GatewayVersion     string                         `json:"gw_version"`
	NameServers        []string                       `json:"nameservers"`
	Netmask            string                         `json:"netmask"`

	Drops            int    `json:"drops"`
	Latency          int    `json:"latency"`
	SpeedTestLastRun int64  `json:"speedtest_lastrun"`
	SpeedTestPing    int    `json:"speedtest_ping"`
	SpeedTestStatus  string `json:"speedtest_status"`

	XPutDown float64 `json:"xput_down"`
	XPutUp   float64 `json:"xput_up"`

	LANIP string `json:"lan_ip"`
}

// SiteHealthResponse contains the site health response data from stat/health
type SiteHealthResponse struct {
	Meta CommonMeta       `json:"meta"`
	Data []SiteHealthData `json:"data"`
}

// SiteHealth queries the site for its health
// site - the site to query
func (c *Client) SiteHealth(site string) (*SiteHealthResponse, error) {
	var resp SiteHealthResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/health", nil, &resp)
	return &resp, err
}
