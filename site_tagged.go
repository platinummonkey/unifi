package unifi

import (
	"net/http"
)

type SiteTaggedMAC map[string]interface{}

type SiteTaggedMACResponse struct {
	Meta CommonMeta      `json:"meta"`
	Data []SiteTaggedMAC `json:"data"`
}

func (c *Client) SiteTaggedMACs(siteID string) (*SiteTaggedMACResponse, error) {
	var resp SiteTaggedMACResponse
	err := c.doSiteRequest(http.MethodGet, siteID, "rest/tag", nil, &resp)
	return &resp, err
}
