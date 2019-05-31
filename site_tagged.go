package unifi

import (
	"net/http"
)

type SiteTaggedMAC map[string]interface{}

type SiteTaggedMACResponse struct {
	Meta CommonMeta      `json:"meta"`
	Data []SiteTaggedMAC `json:"data"`
}

func (c *Client) SiteTaggedMACs(site string) (*SiteTaggedMACResponse, error) {
	var resp SiteTaggedMACResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/tag", nil, &resp)
	return &resp, err
}
