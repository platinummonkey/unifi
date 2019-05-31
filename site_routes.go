package unifi

import (
	"encoding/json"
	"net/http"
)

type SiteRouteNH map[string]interface{}

func (s SiteRouteNH) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(s))
}

type SiteActiveRoutes struct {
	NH  []SiteRouteNH `json:"nh"`
	PFX string        `json:"pfx"`
}

type SiteActiveRoutesResponse struct {
	Meta CommonMeta         `json:"meta"`
	Data []SiteActiveRoutes `json:"data"`
}

func (c *Client) SiteActiveRoutes(site string) (*SiteActiveRoutesResponse, error) {
	var resp SiteActiveRoutesResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/routing", nil, &resp)
	return &resp, err
}

type SiteUserDefinedRoute map[string]interface{}

type SiteUserDefinedRoutesResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data []SiteUserDefinedRoute `json:"data"`
}

func (c *Client) SiteUserDefinedRoutes(site string) (*SiteUserDefinedRoutesResponse, error) {
	var resp SiteUserDefinedRoutesResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/routing", nil, &resp)
	return &resp, err
}
