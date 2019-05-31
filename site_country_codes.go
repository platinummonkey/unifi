package unifi

import (
	"net/http"
)

type SiteCountryCode struct {
	Code interface{} `json:"code"` // sometimes string or int
	Key  string      `json:"key"`
	Name string      `json:"name"`
}

type SiteCountryCodesResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteCountryCode `json:"data"`
}

func (c *Client) SiteCountryCodes(site string) (*SiteCountryCodesResponse, error) {
	var resp SiteCountryCodesResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/ccode", nil, &resp)
	return &resp, err
}
