package unifi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ListBackups will list all auto-backup files
// site - site this device currently registered to
// mac - the device mac
// firmwareURL - the firmware URL
func (c *Client) ListBackups(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "list-backup"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/backup", bytes.NewReader(data), &resp)
	return &resp, err
}

// DeleteBackup will delete a backup on the filesystem
// site - site this device currently registered to
// filename - the backup file to delete
func (c *Client) DeleteBackup(site string, filename string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"cmd":      "delete-backup",
		"filename": filename,
	}
	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/backup", bytes.NewReader(data), &resp)
	return &resp, err
}

// CreateBackup will create a backup to a fixed location on the filesystem.
// site - site this device currently registered to
func (c *Client) CreateBackup(site string) (*GenericResponse, error) {
	data := []byte(`{"cmd": "backup"}`)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPost, site, "cmd/system", bytes.NewReader(data), &resp)
	return &resp, err
}
