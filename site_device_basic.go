package unifi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type SiteDeviceBasic struct {
	Adopted  bool   `json:"adopted"`
	Disabled bool   `json:"disabled"`
	MAC      string `json:"mac"`
	State    int    `json:"state"`
	Type     string `json:"type"`
}

type SiteDeviceBasicResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteDeviceBasic `json:"data"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (c *Client) SiteDevicesBasic(siteID string, typeFilter string) (*SiteDeviceBasicResponse, error) {
	var resp SiteDeviceBasicResponse
	var sendBody io.Reader
	if typeFilter != "" {
		payload := map[string]interface{}{
			"type": typeFilter,
		}
		data, _ := json.Marshal(payload)
		sendBody = bytes.NewReader(data)
	}
	err := c.doSiteRequest(http.MethodGet, siteID, "stat/device-basic", sendBody, &resp)
	return &resp, err
}
