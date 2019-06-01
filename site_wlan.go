package unifi

import (
	"net/http"
)

// SiteWLANConfig is the WLAN configuration
type SiteWLANConfig map[string]interface{}

// SiteWLANConfigResponse contains the WLAN configuration response
type SiteWLANConfigResponse struct {
	Meta CommonMeta       `json:"meta"`
	Data []SiteWLANConfig `json:"data"`
}

// SiteWLANConfigs will query the site for WLAN configurations
// site - the site to query
func (c *Client) SiteWLANConfigs(site string) (*SiteWLANConfigResponse, error) {
	var resp SiteWLANConfigResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/wlanconf", nil, &resp)
	return &resp, err
}

// SiteWLANGroup contains the WLAN group info
type SiteWLANGroup map[string]interface{}

// SiteWLANGroupResponse contains the WLAN group response
type SiteWLANGroupResponse struct {
	Meta CommonMeta      `json:"meta"`
	Data []SiteWLANGroup `json:"data"`
}

// SiteWLANGroups will query the site for WLAN groups
// site - the site to query
func (c *Client) SiteWLANGroups(site string) (*SiteWLANGroupResponse, error) {
	var resp SiteWLANGroupResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/wlangroup", nil, &resp)
	return &resp, err
}
