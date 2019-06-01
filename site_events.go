package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EventSortOrder defines the sort order
type EventSortOrder string

// IsValid returns true if it's a valid sort order.
// there are only a few valid types
func (o EventSortOrder) IsValid() bool {
	switch o {
	case EventSortOrderTimeAscending, EventSortOrderTimeDescending:
		return true
	default:
		return false
	}
}

// Common sort orders
const (
	EventSortOrderTimeDescending EventSortOrder = "-time"
	EventSortOrderTimeAscending  EventSortOrder = "+time"
)

// SiteEventsEvent defines the site event data
type SiteEventsEvent struct {
	ID          string `json:"_id"`
	AP          string `json:"ap"`
	Channel     int    `json:"channel"`
	DatetimeStr string `json:"datetime"`
	HostName    string `json:"hostname"`
	Key         string `json:"key"`
	Message     string `json:"msg"`
	Radio       string `json:"radio"`
	SiteID      string `json:"site_id"`
	SSID        string `json:"ssid"`
	SubSystem   string `json:"subsystem"`
	Time        int64  `json:"time"`
	User        string `json:"user"`
	RadioFrom   string `json:"radio_from"`
	RadioTo     string `json:"radio_to"`
	ChannelFrom string `json:"channel_from"`
	ChannelTo   string `json:"channel_to"`
}

// SiteEventsResponse contains the stat/event site events
type SiteEventsResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteEventsEvent `json:"data"`
}

// SiteEvents will fetch events
// site - site to get IPS/IDS events
// historyHours - number of hours to search in the past
// offset - offset current search if previous request exceeded limit
// limit - limit to number of events to return
// order - how to order the ips/ids events.
func (c *Client) SiteEvents(site string, historyHours int, offset int, limit int, order EventSortOrder) (*SiteEventsResponse, error) {
	if historyHours <= 0 {
		historyHours = 720
	}
	if limit <= 0 {
		limit = 100
	} else if limit > 3000 {
		// there is a default max
		limit = 3000
	}

	if !order.IsValid() {
		return nil, fmt.Errorf("invalid sort order: %s", order)
	}

	payload := map[string]interface{}{
		"_sort":  string(order),
		"within": historyHours,
		"type":   nil,
		"_start": offset,
		"_limit": limit,
	}
	data, _ := json.Marshal(&payload)

	var resp SiteEventsResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/event", bytes.NewReader(data), &resp)
	return &resp, err
}

// SiteIPSEvents will fetch IPS/IDS events
// site - site to get IPS/IDS events
// startTime - start time to search, set to 0 and endTime to 0 for default last 24h behavior
// endTime - end time to search, set to 0 and startTime to 0 for default last 24h behavior
// offset - offset current search if previous request exceeded limit
// limit - limit to number of events to return
// order - how to order the ips/ids events.
func (c *Client) SiteIPSEvents(site string, startTime time.Time, endTime time.Time, offset int, limit int, order EventSortOrder) (*SiteEventsResponse, error) {
	if startTime.IsZero() && endTime.IsZero() {
		endTime = time.Now().UTC()
		startTime = endTime.Add(-24 * time.Hour)
	}
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("end time must come after start time")
	}

	if limit <= 0 {
		limit = 100
	} else if limit > 3000 {
		// there is a default max
		limit = 3000
	}

	if !order.IsValid() {
		return nil, fmt.Errorf("invalid sort order: %s", order)
	}

	payload := map[string]interface{}{
		"_sort":  string(order),
		"type":   nil,
		"start":  startTime.UTC().Unix() * 1000,
		"end":    endTime.UTC().Unix() * 1000,
		"_start": offset,
		"_limit": limit,
	}
	data, _ := json.Marshal(&payload)

	var resp SiteEventsResponse
	err := c.doSiteRequest(http.MethodGet, site, "/stat/ips/event", bytes.NewReader(data), &resp)
	return &resp, err
}
