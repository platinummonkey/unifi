package unifi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// AdoptDevice will adopt a device onto the current site.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) AdoptDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "adopt",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// RestartDevice will restart a device.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) RestartDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "restart",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// ForceProvisionDevice will force-provision an existing device.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) ForceProvisionDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "force-provision",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// PowerCycleDevice will power cycle an existing device.
// site - site this device currently registered to
// mac - the device mac
// portIdx - PoE port to cycle
func (c *Client) PowerCycleDevice(site string, mac string, portIdx int) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":      "power-cycle",
		"mac":      mac,
		"port_idx": portIdx,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// StartSpeedTest will start a speed test.
// site - site this device currently registered to
func (c *Client) StartSpeedTest(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "speedtest"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// SpeedTestStatus will get the current state of a speet test.
// site - site this device currently registered to
func (c *Client) SpeedTestStatus(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "speedtest-status"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// SetLocateDevice will blink a device unit to locate it.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) SetLocateDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "set-locate",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// UnsetLocateDevice will return a blinking device led to normal state.
// site - site this device currently registered to
// mac - the device mac
func (c *Client) UnsetLocateDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "unset-locate",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// UpgradeDevice will trigger a firmware upgrade for the device
// site - site this device currently registered to
// mac - the device mac
func (c *Client) UpgradeDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "upgrade",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// UpgradeExternalDevice will trigger a firmware upgrade for the device with the provided URL location for the firmware.
// site - site this device currently registered to
// mac - the device mac
// firmwareURL - the firmware URL
func (c *Client) UpgradeExternalDevice(site string, mac string, firmwareURL string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "upgrade-external",
		"mac": mac,
		"url": firmwareURL,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// SpectrumScanDevice will trigger a RF scan (AP's only)
// site - site this device currently registered to
// mac - the device mac
// firmwareURL - the firmware URL
func (c *Client) SpectrumScanDevice(site string, mac string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd": "spectrum-scan",
		"mac": mac,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/devmgr", bytes.NewReader(data), &resp)
	return &resp, err
}

// CreateNewUserClientDevice will create a new User/Client device
// mac - client MAC address
// userGroupID - ID value of the UserGroup the new user/client device should belong with.
//               use ListUserGroups to obtain this.
// name - optional name to provide the user/client device
// note - optional note to provide the user/client device
func (c *Client) CreateNewUserClientDevice(site string, mac string, userGroupID string, name string, note string) (*GenericResponse, error) {
	userPayload := map[string]interface{}{
		"mac":          mac,
		"usergroup_id": userGroupID,
	}
	if name != "" {
		userPayload["name"] = name
	}
	if note != "" {
		userPayload["note"] = note
		userPayload["noted"] = true
	}
	payload := map[string]interface{}{
		"objects": []map[string]interface{}{
			{
				"data": userPayload,
			},
		},
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "group/user", bytes.NewReader(data), &resp)
	return &resp, err
}

// SetUserClientDeviceNote will update a note on a user/client device.
// userID - client user ID obtained from SiteDevicesDetailed
// note - optional note to provide the user/client device
//        when note is empty, the existing note for the client-device will be removed
func (c *Client) SetUserClientDeviceNote(site string, userID string, note string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"noted": note != "",
		"note":  note,
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "upd/user", bytes.NewReader(data), &resp)
	return &resp, err
}

// SetUserClientDeviceName will update a name on a user/client device.
// userID - client user ID obtained from SiteDevicesDetailed
// name - optional name to provide the user/client device
//        when note is empty, the existing note for the client-device will be removed
func (c *Client) SetUserClientDeviceName(site string, userID string, name string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"name": name,
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "upd/user", bytes.NewReader(data), &resp)
	return &resp, err
}
