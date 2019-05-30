package unifi

import (
	"net/http"
)

type SiteDetailedSettings struct {
}

type SiteDetailedSettingsResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data []SiteDetailedSettings `json:"data"`
}

func (c *Client) SiteDetailedSettings(siteID string) (*SiteDetailedSettingsResponse, error) {
	var resp SiteDetailedSettingsResponse
	err := c.doSiteRequest(http.MethodGet, siteID, "rest/setting", nil, &resp)
	return &resp, err
}
