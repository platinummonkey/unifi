package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddSite will add a new site to this installation
func (c *Client) AddSite(site string, name string, description string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":  "add-site",
		"name": name,
		"desc": description,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// UpdateSite will update an existing site with a new description.
func (c *Client) UpdateSite(site string, description string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":  "delete-site",
		"desc": description,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// DeleteSite will delete an existing site
func (c *Client) DeleteSite(site string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":  "delete-site",
		"name": site,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// GetSiteAdmins will return the current site admins
func (c *Client) GetSiteAdmins(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "get-admins"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// MoveDevice will move a device from the current site to a new site.
// site - site this device currently registered to
// mac - the device mac
// newSiteID - the new 24 digit site ID to move this device to.
func (c *Client) MoveDevice(site string, mac string, newSiteID string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":     "move-device",
		"mac":     mac,
		"site_id": newSiteID,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// MoveDevice will remove a device from the current site
// site - site this device currently registered to
// mac - the device mac
func (c *Client) DeleteDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "delete-device",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// BlockSTA will block a STA from the current site.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) BlockSTA(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "block-sta",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// UnblockSTA will unblock a STA from the current site.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) UnblockSTA(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "unblock-sta",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// KickSTA will kick a STA from the current site.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) KickSTA(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "kick-sta",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// ForgetSTA will forget a STA from the current site.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) ForgetSTA(site string, macs ...string) (*GenericResponse, error) {
	if len(macs) == 0 {
		return nil, fmt.Errorf("must specify at least one mac")
	}
	payload := map[string]interface{}{
		"cmd": "forget-sta",
		"mac": macs,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}
