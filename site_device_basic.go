package unifi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// SiteDeviceBasic defines the basic device detail data
type SiteDeviceBasic struct {
	Adopted  bool   `json:"adopted"`
	Disabled bool   `json:"disabled"`
	MAC      string `json:"mac"`
	State    int    `json:"state"`
	Type     string `json:"type"`
}

// SiteDeviceBasicResponse contains the stat/device-basic response data
type SiteDeviceBasicResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteDeviceBasic `json:"data"`
}

// SiteDevicesBasic queries the basic device data
// site - the site to query
// typeFilter - the filter to query, if none, then it queries all devices
func (c *Client) SiteDevicesBasic(site string, typeFilter string) (*SiteDeviceBasicResponse, error) {
	var resp SiteDeviceBasicResponse
	var sendBody io.Reader
	if typeFilter != "" {
		payload := map[string]interface{}{
			"type": typeFilter,
		}
		data, _ := json.Marshal(payload)
		sendBody = bytes.NewReader(data)
	}
	err := c.doSiteRequest(http.MethodGet, site, "stat/device-basic", sendBody, &resp)
	return &resp, err
}
