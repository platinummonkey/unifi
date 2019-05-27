package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LoginMeta struct {
	ResponseCode    ResponseCode `json:"rc"`
	ResponseMessage string       `json:"msg"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type LoginResponse struct {
	Meta LoginMeta `json:"meta"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (r *LoginResponse) GetResponseCode() ResponseCode {
	return r.Meta.ResponseCode
}

func (r *LoginResponse) GetResponseMessage() string {
	return r.Meta.ResponseMessage
}

// Login will login the user for making queries
// if remember=true for long-running sessions.
// the API will return HTTP200 for success and a cookie that is your session,
// this method will store this for future commands automatically. Though it is not thread-safe.
func (c *Client) Login(username string, password string, remember bool) error {
	// we do this one manually to acquire cookies
	rememberStr := "false"
	if remember {
		rememberStr = "true"
	}
	u := c.WithPathAndQueryParams("/api/login", "remember", rememberStr)

	auth := map[string]interface{}{
		"username": username,
		"password": password,
	}
	data, _ := json.Marshal(auth)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(data))
	if err != nil {
		return err
	}
	c.SetHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return InvalidResponseBody
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return JSONDecodeError
	}

	if !loginResponse.Meta.ResponseCode.Equal(ResponseCodeOK) {
		return fmt.Errorf("unable to login, response code: %s", loginResponse.Meta.ResponseCode)
	}
	c.authCookies = resp.Cookies()
	c.longRunningSession = remember
	return nil
}

// Logout destroys the sever side session id which will make future attempts with that cookie fail
func (c *Client) Logout() error {
	if !c.longRunningSession {
		// nothing to do, this will be invalid
		return nil
	}
	return c.doRequest(http.MethodGet, "/api/logout", nil, &LoginResponse{})
}

type SelfResponseMeta struct {
	ResponseCode ResponseCode `json:"rc"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type SelfResponseData struct {
	AdminID                   string                 `json:"admin_id"`
	DeviceID                  string                 `json:"device_id"`
	EmailAlertEnabled         bool                   `json:"email_alert_enabled"`
	EmailAlertGroupingDelay   int                    `json:"email_alert_grouping_delay"`
	EmailAlertGroupingEnabled bool                   `json:"email_alert_grouping_enabled"`
	HTMLEmailEnabled          bool                   `json:"html_email_enabled"`
	IsLocal                   bool                   `json:"is_local"`
	IsProfessionalInstaller   bool                   `json:"is_professional_installer"`
	IsSuper                   bool                   `json:"is_super"`
	LastSiteName              string                 `json:"last_site_name,omitempty"`
	Name                      string                 `json:"name"`
	RequiresNewPassword       bool                   `json:"requires_new_password"`
	SuperSitePermissions      []string               `json:"super_site_permissions,omitempty"`
	UISettings                map[string]interface{} `json:"ui_settings"`
}

type SelfResponse struct {
	Data []SelfResponseData `json:"data"` // yes it's an array for some horrible reason.
	Meta SelfResponseMeta   `json:"meta"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (r *SelfResponse) GetResponseCode() ResponseCode {
	return r.Meta.ResponseCode
}

func (r *SelfResponse) GetResponseMessage() string {
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

// Self returns the logged in user.
func (c *Client) Self() (*SelfResponse, error) {
	var selfResponse SelfResponse
	err := c.doRequest(http.MethodGet, "/api/self", nil, &selfResponse)
	return &selfResponse, err
}
