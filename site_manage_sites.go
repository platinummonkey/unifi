package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
)

// AddSite will add a new site to this installation
// site - the current site context
// name - the new site name
// description - the description of the site
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
// site - the site to update
// description - the new site description
func (c *Client) UpdateSite(site string, description string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":  "update-site",
		"desc": description,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// DeleteSite will delete an existing site
// site - the site to delete
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

// SetSiteCountry will set the site's country
// site - the site to update
// siteID - the site's controller id
// configID - the existing country _id configuration - available from SiteDetailedSettings
// country - the country code returned by SiteCountryCodes
func (c *Client) SetSiteCountry(site string, siteID string, configID string, country SiteCountryCode) (*GenericResponse, error) {
	payload := []map[string]interface{}{
		{
			"site_id": siteID,
			"code":    country.Code,
			"key":     "country",
		},
	}
	data, _ := json.Marshal(payload)
	extPath := path.Join("rest/setting/country/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

// SetSiteTimezone will set the site's locale
// site - the site to update
// siteID - the site's controller id
// configID - the existing timezone (locale) _id configuration - available from SiteDetailedSettings
// timezone - the timezone - available from SiteDetailedSettings
func (c *Client) SetSiteTimezone(site string, siteID string, configID string, timezone string) (*GenericResponse, error) {
	payload := []map[string]interface{}{
		{
			"site_id":  siteID,
			"timezone": timezone,
			"key":      "locale",
		},
	}
	data, _ := json.Marshal(payload)
	extPath := path.Join("rest/setting/locale/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

// SetSiteSNMP will set the site's SNMP configuration
// site - the site to update
// siteID - the site's controller id
// configID - the existing SNMP _id configuration - available from SiteDetailedSettings
// community - the SNMP community setting
func (c *Client) SetSiteSNMP(site string, siteID string, configID string, community string) (*GenericResponse, error) {
	payload := []map[string]interface{}{
		{
			"site_id":   siteID,
			"community": community,
			"key":       "snmp",
		},
	}
	data, _ := json.Marshal(payload)
	extPath := path.Join("rest/setting/snmp/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

type SiteManagementConfig struct {
	AdvancedFeatureEnabled *bool `json:"advanced_feature_enabled,omitempty"`
	AlertEnabled           *bool `json:"alert_enabled,omitempty"`
	AutoUpgrade            *bool `json:"auto_upgrade,omitempty"`
	LEDEnabled             *bool `json:"led_enabled,omitempty"`
	UnifiIDPEnabled        *bool `json:"unifi_idp_enabled,omitempty"`
}

// SetSiteManagementConfig will set the site's Management configuration
// site - the site to update
// siteID - the site's controller id
// configID - the existing mgmt _id configuration - available from SiteDetailedSettings
// config - the SiteManagementConfig settings
func (c *Client) SetSiteManagementConfig(site string, siteID string, configID string, config SiteManagementConfig) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"site_id": siteID,
		"key":     "mgmt",
	}
	if config.AdvancedFeatureEnabled != nil {
		payload["advanced_feature_enabled"] = *config.AdvancedFeatureEnabled
	}
	if config.AlertEnabled != nil {
		payload["alert_enabled"] = *config.AlertEnabled
	}
	if config.AutoUpgrade != nil {
		payload["auto_upgrade"] = *config.AutoUpgrade
	}
	if config.LEDEnabled != nil {
		payload["led_enabled"] = *config.LEDEnabled
	}
	if config.UnifiIDPEnabled != nil {
		payload["unifi_idp_enabled"] = *config.UnifiIDPEnabled
	}
	payloads := []map[string]interface{}{payload}
	data, _ := json.Marshal(payloads)
	extPath := path.Join("rest/setting/mgmt/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

type SiteGuessAccessConfig struct {
	Authentication              *string `json:"auth,omitempty"`
	PortalCustomizedTOS         *string `json:"portal_customized_tos,omitempty"`
	PortalCustomizedWelcomeText *string `json:"portal_customized_welcome_text,omitempty"`
	RedirectHTTPS               *bool   `json:"redirect_https,omitempty"`
	RestrictedSubnet1           *string `json:"restricted_subnet_1,omitempty"`
	RestrictedSubnet2           *string `json:"restricted_subnet_2,omitempty"`
	RestrictedSubnet3           *string `json:"restricted_subnet_3,omitempty"`
}

// SetSiteGuestAccessConfig will set the site's guest access configuration
// site - the site to update
// siteID - the site's controller id
// configID - the existing guest_access _id configuration - available from SiteDetailedSettings
// config - the SiteGuessAccessConfig settings
func (c *Client) SetSiteGuestAccessConfig(site string, siteID string, configID string, config SiteGuessAccessConfig) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"site_id": siteID,
		"key":     "guest_access",
	}
	if config.Authentication != nil {
		payload["auth"] = *config.Authentication
	}
	if config.PortalCustomizedTOS != nil {
		payload["portal_customized_tos"] = *config.PortalCustomizedTOS
	}
	if config.PortalCustomizedWelcomeText != nil {
		payload["portal_customized_welcome_text"] = *config.PortalCustomizedWelcomeText
	}
	if config.RedirectHTTPS != nil {
		payload["redirect_https"] = *config.RedirectHTTPS
	}
	if config.RestrictedSubnet1 != nil {
		payload["restricted_subnet_1"] = *config.RestrictedSubnet1
	}
	if config.RestrictedSubnet2 != nil {
		payload["restricted_subnet_2"] = *config.RestrictedSubnet2
	}
	if config.RestrictedSubnet3 != nil {
		payload["restricted_subnet_3"] = *config.RestrictedSubnet3
	}

	payloads := []map[string]interface{}{payload}
	data, _ := json.Marshal(payloads)
	extPath := path.Join("rest/setting/guest_access/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

type SiteNTPConfig struct {
	NTPServer1 *string `json:"ntp_server_1,omitempty"`
	NTPServer2 *string `json:"ntp_server_2,omitempty"`
	NTPServer3 *string `json:"ntp_server_3,omitempty"`
	NTPServer4 *string `json:"ntp_server_4,omitempty"`
}

// SetSiteNTPConfig will set the site's NTP configuration
// site - the site to update
// siteID - the site's controller id
// configID - the existing guest_access _id configuration - available from SiteDetailedSettings
// config - the SiteNTPConfig settings
func (c *Client) SetSiteNTPConfig(site string, siteID string, configID string, config SiteNTPConfig) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"site_id": siteID,
		"key":     "ntp",
	}
	if config.NTPServer1 != nil {
		payload["ntp_server_1"] = *config.NTPServer1
	}
	if config.NTPServer2 != nil {
		payload["ntp_server_2"] = *config.NTPServer2
	}
	if config.NTPServer3 != nil {
		payload["ntp_server_3"] = *config.NTPServer3
	}
	if config.NTPServer4 != nil {
		payload["ntp_server_4"] = *config.NTPServer4
	}

	payloads := []map[string]interface{}{payload}
	data, _ := json.Marshal(payloads)
	extPath := path.Join("rest/setting/connectivity/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

// SetSiteConnectivityConfig will set the site's NTP configuration
// site - the site to update
// siteID - the site's controller id
// configID - the existing guest_access _id configuration - available from SiteDetailedSettings
// uplinkType - the uplink type (e.g. "gateway")
func (c *Client) SetSiteConnectivityConfig(site string, siteID string, configID string, uplinkType string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"site_id":     siteID,
		"key":         "connectivity",
		"uplink_type": strings.TrimSpace(uplinkType),
	}

	payloads := []map[string]interface{}{payload}
	data, _ := json.Marshal(payloads)
	extPath := path.Join("rest/setting/connectivity/", strings.TrimSpace(configID))
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}

// GetSiteAdmins will return the current site admins
// site - the site to query
func (c *Client) GetSiteAdmins(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "get-admins"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// InviteSiteAdmin will invite a new admin for access to the current site
// site - the site to invite the admin to
// name - name to assign to the admin user
// email - email address to assign to the admin user (must be valid to validate)
// disableSSO - set to true to disable SSO capability
// readOnly - set to true to make the admin user read-only
// deviceAdoptPermission - set to true to allow the new admin permissions to adopt devices.
// deviceRestartPermission - set to true to allow the new admin permissions to restart devices.
//
// notes:
//   - after issuing a valid request, an invite will be sent to the email address provided
//   - issuing this command against an existing admin will trigger a "re-invite"
func (c *Client) InviteSiteAdmin(site string, name string, email string, disableSSO bool, readOnly bool, deviceAdoptPermission bool, deviceRestartPermission bool) (*GenericResponse, error) {
	permissions := make([]string, 0)
	if deviceAdoptPermission {
		permissions = append(permissions, "API_DEVICE_ADOPT")
	}
	if deviceRestartPermission {
		permissions = append(permissions, "API_DEVICE_RESTART")
	}
	payload := map[string]interface{}{
		"name":        strings.TrimSpace(name),
		"email":       strings.TrimSpace(email),
		"for_sso":     !disableSSO,
		"cmd":         "invite-admin",
		"role":        "admin",
		"permissions": permissions,
	}
	if readOnly {
		payload["role"] = "readonly"
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// AssignExistingSiteAdmin will assign an existing site admin to the specified site
// site - the site to invite the admin to
// adminID - 24-char string _id of the site admin - from GetSiteAdmins
// readOnly - set to true to make the admin user read-only
// deviceAdoptPermission - set to true to allow the new admin permissions to adopt devices.
// deviceRestartPermission - set to true to allow the new admin permissions to restart devices.
func (c *Client) AssignExistingSiteAdmin(site string, adminID string, readOnly bool, deviceAdoptPermission bool, deviceRestartPermission bool) (*GenericResponse, error) {
	permissions := make([]string, 0)
	if deviceAdoptPermission {
		permissions = append(permissions, "API_DEVICE_ADOPT")
	}
	if deviceRestartPermission {
		permissions = append(permissions, "API_DEVICE_RESTART")
	}
	payload := map[string]interface{}{
		"admin":       strings.TrimSpace(adminID),
		"cmd":         "grant-admin",
		"role":        "admin",
		"permissions": permissions,
	}
	if readOnly {
		payload["role"] = "readonly"
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/sitemgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// RevokeSiteAdmin will revoke a site admin access
// site - the site to invite the admin to
// adminID - 24-char string _id of the site admin - from GetSiteAdmins
func (c *Client) RevokeSiteAdmin(site string, adminID string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"admin": strings.TrimSpace(adminID),
		"cmd":   "revoke-admin",
	}
	data, _ := json.Marshal(payload)

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

// DeleteDevice will remove a device from the current site
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
