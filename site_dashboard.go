package unifi

import (
	"net/http"
)

// ListDashboardMetrics will list dashboard metric objects
// site - the site to query
// scale5Min - if true will return stats based on 5 minute intervals, otherwise defaults to hourly stats.
// note this only works on controllers >= 5.5.x
func (c *Client) ListDashboardMetrics(site string, scale5Min bool) (*GenericResponse, error) {
	var queryParams []string
	if scale5Min {
		queryParams = []string{"scale", "5minutes"}
	}

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/dashboard", nil, &resp, queryParams...)
	return &resp, err
}
