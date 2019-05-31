package unifi

import (
	"bytes"
	"net/http"
)

// ResetDPICounters will reset the site-wide DPI counters
// site - site this device currently registered to
func (c *Client) ResetDPICounters(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "clear-dpi"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/stat", bytes.NewReader(data), &resp)
	return &resp, err
}
