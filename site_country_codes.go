package unifi

import (
	"net/http"
)

// SiteCountryCode defines the country code data
// note: follows `ISO-3166-1 Numeric` standard
type SiteCountryCode struct {
	Code interface{} `json:"code"` // sometimes string or int
	Key  string      `json:"key"`
	Name string      `json:"name"`
}

// SiteCountryCodesResponse defines the country code response
type SiteCountryCodesResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteCountryCode `json:"data"`
}

// SiteCountryCodes lists the site's country codes
// site - the site to query
func (c *Client) SiteCountryCodes(site string) (*SiteCountryCodesResponse, error) {
	var resp SiteCountryCodesResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/ccode", nil, &resp)
	return &resp, err
}
