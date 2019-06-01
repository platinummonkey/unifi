package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ListUserGroups will list all user groups
// site - site to query
func (c *Client) ListUserGroups(site string) (*GenericResponse, error) {
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "list/usergroup", nil, &resp)
	return &resp, err
}

// CreateUserGroup will create a user group
// site - site to modify
// siteID - siteID associated with site
// name - name of the user group
// downloadBandwidth - limit download bandwidth in Kbps (default -1 == unlimited)
// uploadBandwidth - limit upload bandwidth in Kbps (default -1 == unlimited)
func (c *Client) CreateUserGroup(site string, siteID string, name string, downloadBandwidth int, uploadBandwidth int) (*GenericResponse, error) {
	if downloadBandwidth <= 0 {
		downloadBandwidth = -1 // unlimited
	}
	if uploadBandwidth <= 0 {
		uploadBandwidth = -1 // unlimited
	}

	payload := map[string]interface{}{
		"name":              name,
		"qos_rate_max_down": downloadBandwidth,
		"qos_rate_max_up":   uploadBandwidth,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "rest/usergroup", bytes.NewReader(data), &resp)
	return &resp, err
}

// UpdateUserGroup will update an existing user group
// site - site to modify
// siteID - siteID associated with site
// groupID - groupID to modify
// name - name of the user group
// downloadBandwidth - limit download bandwidth in Kbps (default -1 == unlimited)
// uploadBandwidth - limit upload bandwidth in Kbps (default -1 == unlimited)
func (c *Client) UpdateUserGroup(site string, siteID string, groupID string, name string, downloadBandwidth int, uploadBandwidth int) (*GenericResponse, error) {
	if downloadBandwidth <= 0 {
		downloadBandwidth = -1 // unlimited
	}
	if uploadBandwidth <= 0 {
		uploadBandwidth = -1 // unlimited
	}

	payload := map[string]interface{}{
		"_id":               groupID,
		"name":              name,
		"qos_rate_max_down": downloadBandwidth,
		"qos_rate_max_up":   uploadBandwidth,
		"site_id":           siteID,
	}
	data, _ := json.Marshal(payload)

	extPath := fmt.Sprintf("rest/usergroup/%s", strings.TrimSpace(groupID))

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

// DeleteUserGroup will delete an existing user group
// site - site to modify
// groupID - groupID to modify
func (c *Client) DeleteUserGroup(site string, groupID string) (*GenericResponse, error) {
	extPath := fmt.Sprintf("rest/usergroup/%s", strings.TrimSpace(groupID))

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodDelete, site, extPath, nil, &resp)
	return &resp, err
}

// AssignClientUserGroup will assign a user/client device to another group
// site - the site to modify
// clientID - the ID of the user/client device to be modified
// groupID - the ID of the group to assign the user/client device to.
func (c *Client) AssignClientUserGroup(site string, clientID string, groupID string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"usergroup_id": groupID,
	}
	data, _ := json.Marshal(payload)

	extPath := fmt.Sprintf("upd/user/%s", strings.TrimSpace(strings.ToLower(clientID)))

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}
