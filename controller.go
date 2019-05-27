package unifi

import (
	"net/http"
)

type ControllerStatusMeta struct {
	ResponseCode  ResponseCode `json:"rc"`
	ServerVersion string       `json:"server_version"`
	Up            bool         `json:"up"`
	UUID          string       `json:"uuid"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type ControllerStatus struct {
	Meta ControllerStatusMeta `json:"meta"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (r *ControllerStatus) GetResponseCode() ResponseCode {
	return r.Meta.ResponseCode
}

func (r *ControllerStatus) GetResponseMessage() string {
	msg, ok := r.Meta.XXXUnknown["msg"]
	if !ok {
		return ""
	}
	msgStr, ok := msg.(string)
	if !ok {
		return ""
	}
	return msgStr
}

// ControllerStatus returns some very basic server information
// This appears to be the only endpoint that can be reached without an authentication
func (c *Client) ControllerStatus() (*ControllerStatus, error) {
	var status ControllerStatus
	err := c.doRequest(http.MethodGet, "/status", nil, &status)
	return &status, err
}

type SitesResponseMeta struct {
	ResponseCode        ResponseCode `json:"rc"`
	ResponseCodeMessage string       `json:"msg"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type SitesResponseData struct {
	ID                string  `json:"_id"`
	Name              string  `json:"name"`
	Description       string  `json:"desc"`
	Role              string  `json:"role"`
	AttributeHiddenID string  `json:"attr_hidden_id"`
	AttributeNoDelete bool    `json:"attr_no_delete"`
	LocationAccuracy  float64 `json:"location_accuracy"`
	LocationLatitude  float64 `json:"location_lat"`
	LocationLongitude float64 `json:"location_lng"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type SitesResponse struct {
	Data []SitesResponseData `json:"data"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (c *Client) AvailableSites() (*SitesResponse, error) {
	var ret SitesResponse
	err := c.doRequest(http.MethodGet, "/api/self/sites", nil, &ret)
	return &ret, err
}

type SiteAdminsMeta struct {
	ResponseCode    ResponseCode `json:"rc"`
	ResponseMessage string       `json:"msg"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type SiteAdminsResponse struct {
	Meta SiteAdminsMeta `json:"meta"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (r *SiteAdminsResponse) GetResponseCode() ResponseCode {
	return r.Meta.ResponseCode
}

func (r *SiteAdminsResponse) GetResponseMessage() string {
	return r.Meta.ResponseMessage
}

// SiteAdmins returns a list of administrators and permissions for all sites
func (c *Client) SiteAdmins() (*SiteAdminsResponse, error) {
	var resp SiteAdminsResponse
	err := c.doRequest(http.MethodGet, "/api/stat/admin", nil, &resp)
	return &resp, err
}
