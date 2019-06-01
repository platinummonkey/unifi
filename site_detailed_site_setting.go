package unifi

import (
	"net/http"
)

// SiteDetailedSettings contains the detailed site settings
type SiteDetailedSettings map[string]interface{}

// SiteDetailedSettingsResponse contains the rest/setting detailed site settings response
type SiteDetailedSettingsResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data []SiteDetailedSettings `json:"data"`
}

// SiteDetailedSettings queries the site for the detailed settings
// site - the site to query
func (c *Client) SiteDetailedSettings(site string) (*SiteDetailedSettingsResponse, error) {
	var resp SiteDetailedSettingsResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/setting", nil, &resp)
	return &resp, err
}
