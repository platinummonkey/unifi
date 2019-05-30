package unifi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type SuppressionContent struct {
	Alerts    []interface{} `json:"alerts"`
	WhiteList []interface{} `json:"whitelist"`
}

type SiteDeviceDetailedData struct {
	ID                          string             `json:"_id"`
	Key                         string             `json:"key"`
	NTPServer1                  string             `json:"ntp_server_1"`
	NTPServer2                  string             `json:"ntp_server_2"`
	NTPServer3                  string             `json:"ntp_server_3"`
	NTPServer4                  string             `json:"ntp_server_4"`
	SiteID                      string             `json:"site_id"`
	Enabled                     bool               `json:"enabled"`
	BroadcastPing               bool               `json:"broadcast_ping"`
	FTPModule                   bool               `json:"ftp_module"`
	GREModule                   bool               `json:"gre_module"`
	H323Module                  bool               `json:"h323_module"`
	MDNSEnabled                 bool               `json:"mdns_enabled"`
	MSSClamp                    bool               `json:"mss_clamp"`
	OffloadAccounting           bool               `json:"offload_accounting"`
	OffloadL2Blocking           bool               `json:"offload_l2_blocking"`
	OffloadSCH                  bool               `json:"offload_sch"`
	PPTPModule                  bool               `json:"pptp_module"`
	ReceiveRedirects            bool               `json:"receive_redirects"`
	SendRedirects               bool               `json:"send_redirects"`
	SIPModule                   bool               `json:"sip_module"`
	SYNCookies                  bool               `json:"syn_cookies"`
	TFTPModule                  bool               `json:"tftp_module"`
	UPNPEnabled                 bool               `json:"upnp_enabled"`
	UPNPNATPMPEnabled           bool               `json:"upnp_nat_pmp_enabled"`
	UPNPSecureMode              bool               `json:"upnp_secure_mode"`
	Code                        string             `json:"code"`
	TimeZone                    string             `json:"timezone"`
	AdvancedFeatureEnabled      bool               `json:"advanced_feature_enabled"`
	AlertEnabled                bool               `json:"alert_enabled"`
	AutoUpgrade                 bool               `json:"auto_upgrade"`
	LEDEnabled                  bool               `json:"led_enabled"`
	UnifiIDPEnabled             bool               `json:"unifi_idp_enabled"`
	Port                        int                `json:"port"`
	Interval                    int                `json:"interval"`
	Community                   string             `json:"community"`
	UGW3WAN2Enabled             bool               `json:"ugw3_wan2_enabled"`
	UplinkType                  string             `json:"uplink_type"`
	Auth                        string             `json:"auth"`
	PortalCustomizedTOS         string             `json:"portal_customized_tos"`
	PortalCustomizedWelcomeText string             `json:"portal_customized_welcome_text"`
	RedirectHTTPS               bool               `json:"redirect_https"`
	RestrictedSubnet1           string             `json:"restricted_subnet_1"`
	RestrictedSubnet2           string             `json:"restricted_subnet_2"`
	RestrictedSubnet3           string             `json:"restricted_subnet_3"`
	Download                    int                `json:"download"`
	Upload                      int                `json:"upload"`
	EnabledCategories           []string           `json:"enabled_categories"`
	IPSMode                     string             `json:"ips_mode"`
	LastAlertID                 string             `json:"last_alert_id"`
	RestrictIPAddresses         bool               `json:"restrict_ip_addresses"`
	RestrictTOR                 bool               `json:"restrict_tor"`
	Suppression                 SuppressionContent `json:"suppression"`
	UTMToken                    string             `json:"utm_token"`
}

type SiteDeviceDetailedResponse struct {
	Meta CommonMeta               `json:"meta"`
	Data []SiteDeviceDetailedData `json:"data"`
}

func (c *Client) SiteDevicesDetailed(siteID string, filterMACs ...string) (*SiteDeviceDetailedResponse, error) {
	var resp SiteDeviceDetailedResponse
	var sendBody io.Reader
	method := http.MethodGet
	if len(filterMACs) > 0 {
		method = http.MethodPost

		payload := map[string]interface{}{
			"macs": filterMACs,
		}
		data, _ := json.Marshal(payload)
		sendBody = bytes.NewReader(data)
	}
	err := c.doSiteRequest(method, siteID, "stat/device", sendBody, &resp)
	return &resp, err
}
