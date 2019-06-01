package unifi

import (
	"net/http"
	"strconv"
)

// ControllerStatusMeta is the controller status metadata type
type ControllerStatusMeta struct {
	ResponseCode        ResponseCode `json:"rc"`
	ResponseCodeMessage string       `json:"msg"`
	ServerVersion       string       `json:"server_version"`
	Up                  bool         `json:"up"`
	UUID                string       `json:"uuid"`
}

// GetResponseCode returns the response code
func (r *ControllerStatusMeta) GetResponseCode() ResponseCode {
	return r.ResponseCode
}

// GetResponseMessage returns the response message
func (r *ControllerStatusMeta) GetResponseMessage() string {
	return r.ResponseCodeMessage
}

// ControllerStatus is the controller status response
type ControllerStatus struct {
	Meta ControllerStatusMeta     `json:"meta"`
	Data []map[string]interface{} `json:"data"`
}

// ControllerStatus returns some very basic server information
// This appears to be the only endpoint that can be reached without an authentication
func (c *Client) ControllerStatus() (*ControllerStatus, error) {
	var status ControllerStatus
	err := c.doRequest(http.MethodGet, "/status", nil, &status)
	return &status, err
}

// SitesResponseData represents the self/sites response data
type SitesResponseData struct {
	ID                string  `json:"_id"`
	Name              string  `json:"name"`
	Description       string  `json:"desc"`
	Role              string  `json:"role"`
	AttributeHiddenID string  `json:"attr_hidden_id"`
	AttributeNoDelete bool    `json:"attr_no_delete"`
	LocationAccuracy  float64 `json:"location_accuracy"`
	LocationLatitude  float64 `json:"location_lat"`
	LocationLongitude float64 `json:"location_lng"`

	XXXUnknown map[string]interface{} `json:"-"`
}

// SitesResponse represents the self/sites response
type SitesResponse struct {
	Meta CommonMeta          `json:"meta"`
	Data []SitesResponseData `json:"data"`
}

// AvailableSites returns the available sites for the controller.
func (c *Client) AvailableSites() (*SitesResponse, error) {
	var ret SitesResponse
	err := c.doRequest(http.MethodGet, "/api/self/sites", nil, &ret)
	return &ret, err
}

// SitesVerboseGatewaySystemStats contains gateway system stats
type SitesVerboseGatewaySystemStats struct {
	CPUUsage    interface{} `json:"cpu"`    // these come back as strings >.<
	MemoryUsage interface{} `json:"mem"`    // these come back as strings >.<
	Uptime      interface{} `json:"uptime"` // these come back as strings >.<
}

// GetCPUUsage returns the CPU usage
func (s SitesVerboseGatewaySystemStats) GetCPUUsage() float64 {
	switch s.CPUUsage.(type) {
	case string:
		f, err := strconv.ParseFloat(s.CPUUsage.(string), 64)
		if err != nil {
			return 0
		}
		return f
	case float64:
		return s.CPUUsage.(float64)
	}
	return 0
}

// GetMemoryUsage returns the memory usage
func (s SitesVerboseGatewaySystemStats) GetMemoryUsage() float64 {
	switch s.MemoryUsage.(type) {
	case string:
		f, err := strconv.ParseFloat(s.MemoryUsage.(string), 64)
		if err != nil {
			return 0
		}
		return f
	case float64:
		return s.MemoryUsage.(float64)
	}
	return 0
}

// GetUptime returns the uptime
func (s SitesVerboseGatewaySystemStats) GetUptime() float64 {
	switch s.Uptime.(type) {
	case string:
		f, err := strconv.ParseFloat(s.Uptime.(string), 64)
		if err != nil {
			return 0
		}
		return f
	case float64:
		return s.Uptime.(float64)
	}
	return 0
}

// SitesVerboseHealthData is the normalized verbose health data
type SitesVerboseHealthData struct {
	// type 1
	NumberAdopted      int    `json:"num_adopted"`
	NumberAccessPoints int    `json:"num_ap"`
	NumberDisabled     int    `json:"num_disabled"`
	NumberDisconnected int    `json:"num_disconnected"`
	NumberGuest        int    `json:"num_guest"`
	NumberIoTDevices   int    `json:"num_iot"`
	NumberPending      int    `json:"num_pending"`
	NumberUsers        int    `json:"num_user"`
	RxBytesR           int64  `json:"rx_bytes-r"`
	Status             string `json:"status"`
	SubSystem          string `json:"subsystem"`
	TxBytesR           int64  `json:"tx_bytes-r"`

	// type 2
	Gateways           []string                       `json:"gateways"`
	GatewayMAC         string                         `json:"gw_mac"`
	GatewayName        string                         `json:"gw_name"`
	GatewaySystemStats SitesVerboseGatewaySystemStats `json:"gw_system-stats"`
	GatewayVersion     string                         `json:"gw_version"`
	NameServers        []string                       `json:"nameservers"`
	NetMask            string                         `json:"netmask"`
	NumberGateways     int                            `json:"num_gw"`
	NumberSTA          int                            `json:"num_sta"`
	WANIP              string                         `json:"wan_ip"`

	// Type 3
	Drops            int64   `json:"drops"`
	Latency          int64   `json:"latency"`
	SpeedTestLastRun float64 `json:"speedtest_lastrun"`
	SpeedTestPing    int64   `json:"speedtest_ping"`
	SpeedTestStatus  string  `json:"speedtest_status"`
	XPutDown         float64 `json:"xput_down"`
	XPutUp           float64 `json:"xput_up"`

	// Type 4
	LANIP    string `json:"lan_ip"`
	NumberSW int    `json:"num_sw"`
}

// SitesVerboseResponseData contains normalized health data
type SitesVerboseResponseData struct {
	ID                string  `json:"_id"`
	Name              string  `json:"name"`
	Description       string  `json:"desc"`
	Role              string  `json:"role"`
	AttributeHiddenID string  `json:"attr_hidden_id"`
	AttributeNoDelete bool    `json:"attr_no_delete"`
	LocationAccuracy  float64 `json:"location_accuracy"`
	LocationLatitude  float64 `json:"location_lat"`
	LocationLongitude float64 `json:"location_lng"`
	NumberNewAlarms   int     `json:"num_new_alarms"`

	Health []SitesVerboseHealthData `json:"health"`
}

// SitesVerboseResponse is the verbose response for stat/sites
type SitesVerboseResponse struct {
	Meta CommonMeta                 `json:"meta"`
	Data []SitesVerboseResponseData `json:"data"`
}

// AvailableSitesVerbose returns the available sites with verbose health data
func (c *Client) AvailableSitesVerbose() (*SitesVerboseResponse, error) {
	var ret SitesVerboseResponse
	err := c.doRequest(http.MethodGet, "/api/stat/sites", nil, &ret)
	return &ret, err
}

// SiteAdminsResponse returns the stat/admin response data
type SiteAdminsResponse struct {
	Meta CommonMeta               `json:"meta"`
	Data []map[string]interface{} `json:"data"`
}

// SiteAdmins returns a list of administrators and permissions for all sites
func (c *Client) SiteAdmins() (*SiteAdminsResponse, error) {
	var resp SiteAdminsResponse
	err := c.doRequest(http.MethodGet, "/api/stat/admin", nil, &resp)
	return &resp, err
}
