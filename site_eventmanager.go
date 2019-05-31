package unifi

import (
	"bytes"
	"net/http"
)

// ArchiveAllAlarms will archive all alarms
func (c *Client) ArchiveAllAlarms(site string) error {
	data := []byte(`{"cmd": "archive-all-alarms"}`)
	return c.doSiteRequest(http.MethodPost, site, "cmd/evtmgt", bytes.NewReader(data), nil)
}
