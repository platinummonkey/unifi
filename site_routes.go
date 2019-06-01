package unifi

import (
	"net/http"
)

// SiteRouteNH defines the site route information
type SiteRouteNH map[string]interface{}

// SiteActiveRoutes defines the active routes definitions
type SiteActiveRoutes struct {
	NH  []SiteRouteNH `json:"nh"`
	PFX string        `json:"pfx"`
}

// SiteActiveRoutesResponse contains the active routes definition response
type SiteActiveRoutesResponse struct {
	Meta CommonMeta         `json:"meta"`
	Data []SiteActiveRoutes `json:"data"`
}

// SiteActiveRoutes lists active routes for the site
// site - the site to query
func (c *Client) SiteActiveRoutes(site string) (*SiteActiveRoutesResponse, error) {
	var resp SiteActiveRoutesResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/routing", nil, &resp)
	return &resp, err
}

// SiteUserDefinedRoute is a user defined route
type SiteUserDefinedRoute map[string]interface{}

// SiteUserDefinedRoutesResponse contains the user defined routes response
type SiteUserDefinedRoutesResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data []SiteUserDefinedRoute `json:"data"`
}

// SiteUserDefinedRoutes queries the user defines routes
// site - the site to query
func (c *Client) SiteUserDefinedRoutes(site string) (*SiteUserDefinedRoutesResponse, error) {
	var resp SiteUserDefinedRoutesResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/routing", nil, &resp)
	return &resp, err
}
