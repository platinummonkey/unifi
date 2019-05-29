package unifi

import (
	"net/http"
)

type SiteSysInfo struct {
	AutoBackup                               bool        `json:"autobackup"`
	Build                                    string      `json:"build"`
	CloudKeySDCardMounted                    bool        `json:"cloudkey_sdcard_mounted"`
	CloudKeyStorageLimitations               bool        `json:"cloudkey_storage_limitations"`
	CloudKeyUpdateAvailable                  bool        `json:"cloudkey_update_available"`
	CloudKeyUpdateLastChecked                interface{} `json:"cloudkey_update_last_checked"` // sometimes string or int
	CloudKeyUsingFakeSystemProperties        bool        `json:"cloudkey_using_fake_system_properties"`
	CloudKeyVersion                          string      `json:"cloudkey_version"`
	DataRetentionDays                        int         `json:"data_retention_days"`
	DataRetentionTimeInHoursFor5MinutesScale int         `json:"data_retention_time_in_hours_for_5minutes_scale"`
	DataRetentionTimeInHoursForDailyScale    int         `json:"data_retention_time_in_hours_for_daily_scale"`
	DataRetentionTimeInHoursForHourlyScale   int         `json:"data_retention_time_in_hours_for_hourly_scale"`
	DataRetentionTimeInHoursForMonthlyScale  int         `json:"data_retention_time_in_hours_for_monthly_scale"`
	DataRetentionTimeInHoursForOthers        int         `json:"data_retention_time_in_hours_for_others"`
	DebugDevice                              string      `json:"debug_device"`
	DebugManagement                          string      `json:"debug_mgmt"`
	DebugSDN                                 string      `json:"debug_sdn"`
	DebugSystem                              string      `json:"debug_system"`
	EOLPendingDeviceCount                    int         `json:"eol_pending_device_count"`
	FacebookWIFIRegistered                   bool        `json:"facebook_wifi_registered"`
	GoogleMapsAPIKey                         string      `json:"google_maps_api_key"`
	HostName                                 string      `json:"hostname"`
	HTTPSPort                                int         `json:"https_port"`
	ImageMapsUseGoogleEngine                 bool        `json:"image_maps_use_google_engine"`
	InformPort                               int         `json:"inform_port"`
	IPAddresses                              []string    `json:"ip_addrs"`
	LiveChat                                 string      `json:"live_chat"`
	Name                                     string      `json:"name"`
	OverrideInformHost                       bool        `json:"override_inform_host"`
	PackageUpdateAvailable                   bool        `json:"package_update_available"`
	PackageUpdateLastChecked                 interface{} `json:"package_update_last_checked"` // sometimes string or int
	PackageVersion                           string      `json:"package_version"`
	PreviousVersion                          string      `json:"previous_version"`
	RadiusDisconnectRunning                  bool        `json:"radius_disconnect_running"`
	StoreEnabled                             string      `json:"store_enabled"`
	TimeZone                                 string      `json:"timezone"`
	UnifiGoEnabled                           bool        `json:"unifi_go_enabled"`
	UnsupportedDeviceCount                   int         `json:"unsupported_device_count"`
	UpdateAvailable                          bool        `json:"update_available"`
	UpdateDownloaded                         bool        `json:"update_downloaded"`
	Version                                  string      `json:"version"`

	XXXUnknown map[string]interface{} `json:"-"`
}

type SiteSysInfoResponse struct {
	Meta CommonMeta    `json:"meta"`
	Data []SiteSysInfo `json:"data"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (c *Client) SiteSysInfo(siteID string) (*SiteSysInfoResponse, error) {
	var resp SiteSysInfoResponse
	err := c.doSiteRequest(http.MethodGet, siteID, "stat/sysinfo", nil, &resp)
	return &resp, err
}
