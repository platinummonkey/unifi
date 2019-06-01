package unifi

import (
	"net/http"
)

// SiteTaggedMAC a tagged MAC device
type SiteTaggedMAC map[string]interface{}

// SiteTaggedMACResponse contains tagged MAC device info response
type SiteTaggedMACResponse struct {
	Meta CommonMeta      `json:"meta"`
	Data []SiteTaggedMAC `json:"data"`
}

// SiteTaggedMACs will query the site for tagged MACs
// site - the site to query
func (c *Client) SiteTaggedMACs(site string) (*SiteTaggedMACResponse, error) {
	var resp SiteTaggedMACResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/tag", nil, &resp)
	return &resp, err
}
