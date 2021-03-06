package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SiteActiveClient defines an active client device
type SiteActiveClient struct {
	ID            string `json:"_id"`
	IsGuestByUAP  bool   `json:"_is_guest_by_uap"`
	LastSeenByUAP int64  `json:"_last_seen_by_uap"`
	UptimeByUAP   int64  `json:"_uptime_by_uap"`

	IsGuestByUGW  bool  `json:"_is_guest_by_ugw"`
	LastSeenByUGW int64 `json:"_last_seen_by_ugw"`
	UptimeByUGW   int64 `json:"_uptime_by_ugw"`

	Anomalies             int    `json:"anomalies"`
	AccessPointMAC        string `json:"ap_mac"`
	AssociationTime       int64  `json:"assoc_time"`
	Authorized            bool   `json:"authorized"`
	BSSID                 string `json:"bssid"`
	BytesR                int64  `json:"bytes-r"`
	CCQ                   int    `json:"ccq"`
	Channel               int    `json:"channel"`
	DHCPEndTime           int    `json:"dhcpend_time"`
	ESSID                 string `json:"essid"`
	FirstSeen             int64  `json:"first_seen"`
	GatewayMAC            string `json:"gw_mac"`
	HostName              string `json:"hostname"`
	IdleTime              int    `json:"idletime"`
	IP                    string `json:"ip"`
	Is11R                 bool   `json:"is_11r"`
	IsGuest               bool   `json:"is_guest"`
	IsWired               bool   `json:"is_wired"`
	LastSeen              int64  `json:"last_seen"`
	LatestAssociationTime int64  `json:"latest_assoc_time"`
	MAC                   string `json:"mac"`
	Network               string `json:"network"`
	NetworkID             string `json:"network_id"`
	Noise                 int    `json:"noise"`
	OUI                   string `json:"oui"`
	PowerSaveEnabled      bool   `json:"powersave_enabled"`
	QOSPolicyApplied      bool   `json:"qos_policy_applied"`
	Radio                 string `json:"radio"`
	RadioName             string `json:"radio_name"`
	RadioProto            string `json:"radio_proto"`
	RSSI                  int    `json:"rssi"`
	RXBytes               int64  `json:"rx_bytes"`
	RXBytesR              int64  `json:"rx_bytes-r"`
	RXPackets             int64  `json:"rx_packets"`
	RXRate                int64  `json:"rx_rate"`
	Satisfaction          int    `json:"satisfaction"`
	Signal                int    `json:"signal"`
	SiteID                string `json:"site_id"`
	TXBytes               int64  `json:"tx_bytes"`
	TXBytesR              int64  `json:"tx_bytes-r"`
	TXPackets             int64  `json:"tx_packets"`
	TXPower               int    `json:"tx_power"`
	TXRate                int64  `json:"tx_rate"`
	Uptime                int64  `json:"uptime"`
	UserID                string `json:"user_id"`
	VLAN                  int    `json:"vlan"`
}

// SiteActiveClientsResponse contains the active clients response
type SiteActiveClientsResponse struct {
	Meta CommonMeta         `json:"meta"`
	Data []SiteActiveClient `json:"data"`
}

// SiteActiveClients will list active clients
// site - the site to query
// filterMac - filter to a specific mac, if zero-value, then no filter is applied
func (c *Client) SiteActiveClients(site string, filterMac string) (*SiteActiveClientsResponse, error) {
	extPath := "stat/sta"
	if filterMac != "" {
		extPath = extPath + "/" + strings.ToLower(filterMac)
	}

	var resp SiteActiveClientsResponse
	err := c.doSiteRequest(http.MethodGet, site, extPath, nil, &resp)
	return &resp, err
}

// ClientDetails gets the details for a single client
// site - the site to query
// mac - the client mac to query
func (c *Client) ClientDetails(site string, mac string) (*GenericResponse, error) {
	if mac == "" {
		return nil, fmt.Errorf("must specify a client MAC")
	}

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, fmt.Sprintf("stat/user/%s", strings.ToLower(mac)), nil, &resp)
	return &resp, err
}

// UpdateClientFixedIP will update a clients fixedIP
// site - the site to modify
// clientID - the ID of the user/client device to be modified
// useFixedIP - true to set a fixedIP, false to unset
// networkID - if useFixedIP set this to the specified value
// fixedIP - if userFixedIP set this to the fixed IP specified
func (c *Client) UpdateClientFixedIP(site string, clientID string, useFixedIP bool, networkID *string, fixedIP *string) (*GenericResponse, error) {
	payload := map[string]interface{}{
		"_id":         strings.TrimSpace(strings.ToLower(clientID)),
		"use_fixedip": useFixedIP,
	}
	if useFixedIP {
		if networkID != nil {
			payload["network_id"] = *networkID
		}
		if fixedIP != nil {
			payload["fixed_ip"] = *fixedIP
		}
	}

	data, _ := json.Marshal(payload)

	extPath := fmt.Sprintf("rest/user/%s", strings.TrimSpace(strings.ToLower(clientID)))

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodPut, site, extPath, bytes.NewReader(data), &resp)
	return &resp, err
}
