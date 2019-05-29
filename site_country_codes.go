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

	XXXUnknown map[string]interface{} `json:"-"`
}

func (c *Client) SiteCountryCodes(siteID string) (*SiteCountryCodesResponse, error) {
	var resp SiteCountryCodesResponse
	err := c.doSiteRequest(http.MethodGet, siteID, "stat/ccode", nil, &resp)
	return &resp, err
}
