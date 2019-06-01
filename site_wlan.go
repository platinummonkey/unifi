package unifi

import (
	"net/http"
)

type SiteWLANConfig map[string]interface{}

type SiteWLANConfigResponse struct {
	Meta CommonMeta       `json:"meta"`
	Data []SiteWLANConfig `json:"data"`
}

func (c *Client) SiteWLANConfigs(site string) (*SiteWLANConfigResponse, error) {
	var resp SiteWLANConfigResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/wlanconf", nil, &resp)
	return &resp, err
}

type SiteWLANGroup map[string]interface{}

type SiteWLANGroupResponse struct {
	Meta CommonMeta      `json:"meta"`
	Data []SiteWLANGroup `json:"data"`
}

func (c *Client) SiteWLANGroups(site string) (*SiteWLANGroupResponse, error) {
	var resp SiteWLANGroupResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/wlangroup", nil, &resp)
	return &resp, err
}
