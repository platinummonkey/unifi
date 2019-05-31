package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EventSortOrder string

func (o EventSortOrder) IsValid() bool {
	switch o {
	case EventSortOrderTimeAscending, EventSortOrderTimeDescending:
		return true
	default:
		return false
	}
}

const (
	EventSortOrderTimeDescending EventSortOrder = "-time"
	EventSortOrderTimeAscending  EventSortOrder = "+time"
)

var InvalidSortOrderError = fmt.Errorf("invalid sort order")

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

type SiteEventsResponse struct {
	Meta CommonMeta        `json:"meta"`
	Data []SiteEventsEvent `json:"data"`
}

// type SiteEventsResponse map[string]interface{}

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

	var resp SiteEventsResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/event", bytes.NewReader(data), &resp)
	return &resp, err
}
