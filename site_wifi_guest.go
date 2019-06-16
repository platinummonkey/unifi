package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// WifiGuestConfig is the guest wifi configuration
type WifiGuestConfig struct {
	UploadSpeed    int    // upload speed in Kbps
	DownloadSpeed  int    // download speed in Kbps
	TransferLimit  int    // data transfer limits in MB
	AccessPointMac string // Access Point MAC address to which the client is connected, should result in faster authorization
}

// AuthorizeWiFiGuest will authorize a WiFi guest on the network
// site - site to allow the guest
// mac - client mac to authorize
// duration - time for wifi authorization, if <=0 , then it will default to 1hr
// wifiGuestConfig - optional parameters to limit the client
func (c *Client) AuthorizeWiFiGuest(site string, mac string, duration time.Duration, wifiGuestConfig *WifiGuestConfig) (*GenericResponse, error) {
	if duration.Minutes() <= 0 {
		duration = time.Hour * 1
	}

	payload := map[string]interface{}{
		"cmd":     "authorize-guest",
		"mac":     strings.ToLower(mac),
		"minutes": int(duration.Minutes()),
	}
	// optionally each parameter if non-zero
	if wifiGuestConfig != nil {
		if wifiGuestConfig.UploadSpeed >= 0 {
			payload["up"] = wifiGuestConfig.UploadSpeed
		}

		if wifiGuestConfig.DownloadSpeed >= 0 {
			payload["down"] = wifiGuestConfig.DownloadSpeed
		}

		if wifiGuestConfig.TransferLimit >= 0 {
			payload["bytes"] = wifiGuestConfig.TransferLimit
		}

		if wifiGuestConfig.AccessPointMac != "" {
			payload["ap_mac"] = wifiGuestConfig.AccessPointMac
		}
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/stamgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// UnAuthorizeWiFiGuest will unauthorize a WiFi guest
// site - site to allow the guest
// mac - client mac to unauthorize
func (c *Client) UnAuthorizeWiFiGuest(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "unauthorize-guest",
		"mac": strings.ToLower(mac),
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/stamgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListWiFiGuests will list guest devices with valid access
// site - site to query
// withinHours - time frame in hours to list guest devices, default value if zero is 24 hours
func (c *Client) ListWiFiGuests(site string, withinHours int) (*GenericResponse, error) {
	if withinHours <= 0 {
		withinHours = 24
	}

	payload := map[string]interface{}{
		"within": withinHours,
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/guest", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListWiFiGuestVouchers will list wifi guest vouchers
// site - the site to query
// createdTime - the create time of the voucher, if zero-value, then it will return all
func (c *Client) ListWiFiGuestVouchers(site string, createTime time.Time) (*GenericResponse, error) {

	payload := map[string]interface{}{}
	if !createTime.IsZero() {
		payload["create_time"] = createTime.UTC().Unix()
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/voucher", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListWiFiGuestPayments will list wifi guest payments
// site - the site to query
// withinHours - number of hours to search for history, if zero, then use default 24 hours
func (c *Client) ListWiFiGuestPayments(site string, withinHours int) (*GenericResponse, error) {
	if withinHours <= 0 {
		withinHours = 24
	}
	extPath := fmt.Sprintf("stat/payment?within=%d", withinHours)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, extPath, nil, &resp)
	return &resp, err
}

// CreateWifiGuestOperator will create a wifi guest operator
// site - the site to create a new wifi guest operator
// name - the name the new wifi guest operator
// password - the clear text password for the wifi guest operator
// note - optional note to attach to the wifi guest operator
func (c *Client) CreateWifiGuestOperator(site string, name string, password string, note string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"name":     strings.TrimSpace(name),
		"password": password,
	}
	note = strings.TrimSpace(note)
	if note != "" {
		payload["note"] = note
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "rest/hotspotop", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListWiFiGuestOperators will list wifi guest operators
// site - the site to query
func (c *Client) ListWiFiGuestOperators(site string) (*GenericResponse, error) {
	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "rest/hotspotop", nil, &resp)
	return &resp, err
}

// VoucherConfig defines a voucher configuration
type VoucherConfig struct {
	MinutesValid       uint    // minutes the voucher is valid after activation (expiration time)
	Count              *uint   // number of vouchers to create, default value is 1
	Quota              uint    // single-use or multi-use vouchers, value `0` is for multi-use, `1` is for single-use, anything else is for multi-use <N> times
	Note               *string // note text to add to voucher when printing
	UploadSpeedLimit   *uint   // upload speed limit in kbps
	DownloadSpeedLimit *uint   // download speed limit in kbps
	DataTransferLimit  *uint   // data transfer limit in MB
}

// CreateWifiGuestVoucher will create a wifi guest voucher
// site - the site to create a new wifi guest voucher
// cfg - voucher creation config
func (c *Client) CreateWifiGuestVoucher(site string, cfg VoucherConfig) (*GenericResponse, error) {
	count := uint(1)
	if cfg.Count != nil {
		count = *cfg.Count
	}

	payload := map[string]interface{}{
		"cmd":    "create-voucher",
		"expire": cfg.MinutesValid,
		"n":      count,
		"quota":  cfg.Quota,
	}
	if cfg.Note != nil {
		note := strings.TrimSpace(*cfg.Note)
		payload["note"] = note
	}
	if cfg.UploadSpeedLimit != nil {
		payload["up"] = *cfg.UploadSpeedLimit
	}
	if cfg.DownloadSpeedLimit != nil {
		payload["down"] = *cfg.DownloadSpeedLimit
	}
	if cfg.DataTransferLimit != nil {
		payload["bytes"] = *cfg.DataTransferLimit
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/hotspot", bytes.NewReader(data), &resp)
	return &resp, err
}

// RevokeWifiGuestVoucher will revoke a guest wifi voucher
// site - the site to create a revoke wifi guest voucher
// voucherID - the voucher _id to revoke
func (c *Client) RevokeWifiGuestVoucher(site string, voucherID string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "delete-voucher",
		"_id": strings.TrimSpace(voucherID),
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/hotspot", bytes.NewReader(data), &resp)
	return &resp, err
}

// ExtendWifiGuestValidity will extend a guest wifi client
// site - the site to create a revoke wifi guest voucher
// guestID - the guest _id to extend validity
func (c *Client) ExtendWifiGuestValidity(site string, guestID string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "extend",
		"_id": strings.TrimSpace(guestID),
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/hotspot", bytes.NewReader(data), &resp)
	return &resp, err
}
