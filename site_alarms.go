package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SiteAlarmsAlarm is an alarm event
type SiteAlarmsAlarm struct {
	ID                    string      `json:"_id"`
	Archived              bool        `json:"archived"`
	DatetimeStr           string      `json:"datetime"`
	Gateway               string      `json:"gw"`
	GatewayName           string      `json:"gw_name"`
	Key                   string      `json:"key"`
	Message               string      `json:"msg"`
	SiteID                string      `json:"site_id"`
	SubSystem             string      `json:"subsystem"`
	Time                  int64       `json:"time"`
	Timestamp             int64       `json:"timestamp"`
	CatName               string      `json:"catname"`
	DestinationIP         string      `json:"dest_ip"`
	DestinationPort       int         `json:"dest_port"`
	DestinationMAC        string      `json:"dst_mac"`
	DestinationIPASN      string      `json:"dstipASN"`
	DestinationIPCountry  interface{} `json:"dstipCountry"` // seems to be false sometimes
	DestinationIPGeo      GeoCodeData `json:"dstipGeo"`     // seems to be false sometimes
	EventType             string      `json:"event_type"`
	FlowID                int64       `json:"flow_id"`
	Host                  string      `json:"host"`
	InterfaceIn           string      `json:"in_iface"`
	InnerAlertAction      string      `json:"inner_alert_action"`
	InnerAlertCategory    string      `json:"inner_alert_category"`
	InnerAlertGID         int         `json:"inner_alert_gid"`
	InnerAlertRevision    int         `json:"inner_alert_rev"`
	InnerAlertSeverity    int         `json:"inner_alert_severity"`
	InnerAlertSignature   string      `json:"inner_alert_signature"`
	InnerAlertSignatureID int         `json:"inner_alert_signature_id"`
	From                  string      `json:"from"`
	To                    string      `json:"to"`
	Protocol              string      `json:"protocol"`
	Proto                 string      `json:"proto"`
	SourceIP              string      `json:"src_ip"`
	SourceMAC             string      `json:"src_mac"`
	SourcePort            int         `json:"src_port"`
	SourceIPASN           string      `json:"srcipASN"`
	SourceIPCountry       interface{} `json:"srcipCountry"`
	SourceIPGeo           GeoCodeData `json:"srcipGeo"`
	UniqueAlertID         string      `json:"unique_alertid"`
	USGIP                 string      `json:"usgip"`
	USGMAC                string      `json:"usg_mac"`
	USGPort               int         `json:"usg_port"`
	USGIPASN              string      `json:"usgipASN"`
	USGIPCountry          interface{} `json:"usgipCountry"`
	USGIPGeo              GeoCodeData `json:"usgipGeo"`
}

// SiteAlarmsResponse contains the stat/alarms alarm events response
type SiteAlarmsResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteAlarmsAlarm `json:"data"`
}

// SiteAlarms returns the alarm events for the site
// site - site to query
// historyHours - number of hours of history to return, defaults to 24 hours
// offset - offset of current query, default is 0
// limit - limit the max amount of events returned, defaults to 3000 if zero-value
// order - defined the sort order of the alarm events
// archived - query archived (when true) or unarchived (default) alarm events
func (c *Client) SiteAlarms(site string, historyHours int, offset int, limit int, order EventSortOrder, archived bool) (*SiteAlarmsResponse, error) {
	if historyHours <= 0 {
		historyHours = 24
	}
	if limit <= 0 {
		limit = 100
	} else if limit > 3000 {
		// there is a default max
		limit = 3000
	}

	if !order.IsValid() {
		return nil, InvalidSortOrderError
	}

	payload := map[string]interface{}{
		"_sort":  string(order),
		"within": historyHours,
		"type":   nil,
		"_start": offset,
		"_limit": limit,
	}
	data, _ := json.Marshal(&payload)

	var resp SiteAlarmsResponse
	err := c.doSiteRequest(http.MethodGet, site, fmt.Sprintf("stat/alarm?archived=%t", archived), bytes.NewReader(data), &resp)
	return &resp, err
}

// SiteAlarmsCountResponse returns the stat/alarm/cnt/ response data
type SiteAlarmsCountResponse struct {
	Meta CommonMeta               `json:"meta"`
	Data []map[string]interface{} `json:"data"`
}

// SiteAlarmsCount returns the count of site alarms
// site - site to query
// historyHours - number of hours of history to return, defaults to 24 hours
// archived - query archived (when true) or unarchived (default) alarm events
func (c *Client) SiteAlarmsCount(site string, historyHours int, archived bool) (*SiteAlarmsCountResponse, error) {
	if historyHours <= 0 {
		historyHours = 720
	}

	payload := map[string]interface{}{
		"within": historyHours,
		"type":   nil,
	}
	data, _ := json.Marshal(&payload)

	var resp SiteAlarmsCountResponse
	err := c.doSiteRequest(http.MethodGet, site, fmt.Sprintf("stat/alarm/cnt/?archived=%t", archived), bytes.NewReader(data), &resp)
	return &resp, err
}
