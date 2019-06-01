package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SessionType defines the session type
type SessionType string

// Common session type requests
const (
	SessionTypeAll   SessionType = "all"
	SessionTypeGuest SessionType = "guest"
	SessionTypeUser  SessionType = "user"
)

// IsValid returns true if it's a valid session type.
// there are only a few valid types
func (t SessionType) IsValid() bool {
	switch t {
	case SessionTypeAll, SessionTypeGuest, SessionTypeUser:
		return true
	default:
		return false
	}
}

// SiteSessionOrder defines the sort order for the site sessions
type SiteSessionOrder string

// The supported site session sort orders
const (
	SiteSessionSortOrderTimeDescending SiteSessionOrder = "-assoc_time"
	SiteSessionSortOrderTimeAscending  SiteSessionOrder = "+assoc_time"
)

// IsValid returns true if it's a valid sort order.
// there are only a few valid types
func (o SiteSessionOrder) IsValid() bool {
	switch o {
	case SiteSessionSortOrderTimeAscending, SiteSessionSortOrderTimeDescending:
		return true
	default:
		return false
	}
}

// ListLoginSessions will show all login sessions
// site - site to query
// sessionType - the type of session to query
// startTime - start time to query, set to 0 and endTime to 0 to get default last 1 hour behavior
// endTime - end time to query, set to 0 and startTime to 0 to get default last 1 hour behavior
// mac - mac to filter on, set to `""` for no filtering.
func (c *Client) ListLoginSessions(site string, sessionType SessionType, startTime time.Time, endTime time.Time, mac string) (*GenericResponse, error) {
	if startTime.IsZero() && endTime.IsZero() {
		endTime := time.Now().UTC()
		startTime = endTime.Add(-1 * time.Hour)
	}
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("endTime must be after starTime")
	}

	if !sessionType.IsValid() {
		return nil, fmt.Errorf("invalid sesisonType: %s", sessionType)
	}

	payload := map[string]interface{}{
		"type":  string(sessionType),
		"start": startTime.UTC().Unix() * 1000,
		"end":   endTime.UTC().Unix() * 1000,
	}
	if mac != "" {
		payload["mac"] = strings.ToLower(mac)
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/session", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListLatestSessions will show the latest login sessions
// site - site to query
// mac - mac to filter on
// order - how to order the session events
// offset - offset current request, default to 0 if zero-value
// limit - limit the number of returned sessions, default to 100 if zero-value
func (c *Client) ListLatestSessions(site string, mac string, order SiteSessionOrder, offset int, limit int) (*GenericResponse, error) {
	if mac == "" {
		return nil, fmt.Errorf("must specifiy a client device MAC")
	}

	if offset <= 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 100
	}

	if !order.IsValid() {
		return nil, fmt.Errorf("invalid session order: %s", order)
	}

	payload := map[string]interface{}{
		"mac":    strings.ToLower(mac),
		"_limit": limit,
		"_sort":  string(order),
		"_start": offset,
	}
	if mac != "" {
		payload["mac"] = strings.ToLower(mac)
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/session", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListAuthorizations will list all authorizations
// site - site to query
// startTime - start time to query, set to 0 and endTime to 0 to get default last 1 hour behavior
// endTime - end time to query, set to 0 and startTime to 0 to get default last 1 hour behavior
func (c *Client) ListAuthorizations(site string, startTime time.Time, endTime time.Time) (*GenericResponse, error) {
	if startTime.IsZero() && endTime.IsZero() {
		endTime := time.Now().UTC()
		startTime = endTime.Add(-1 * time.Hour)
	}
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("endTime must be after starTime")
	}

	payload := map[string]interface{}{
		"start": startTime.UTC().Unix() * 1000,
		"end":   endTime.UTC().Unix() * 1000,
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/authorization", bytes.NewReader(data), &resp)
	return &resp, err
}

// ListAllUsers will show the clients ever connected to the site
// site - site to query
// withinHours - hours to go back, default to 24 if zero-value
// offset - offset current request, default to 0 if zero-value
// limit - limit the number of returned sessions, default to 100 if zero-value
//
// note: withinHours filters clients that were connected within the period
//       the returned stats per client are all-time totals, irrespective of withinHours
func (c *Client) ListAllUsers(site string, withinHours int, offset int, limit int) (*GenericResponse, error) {
	if withinHours <= 0 {
		withinHours = 24
	}

	if offset <= 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 1000
	}

	payload := map[string]interface{}{
		"type":   "all",
		"conn":   "all",
		"within": withinHours,
		"_start": offset,
		"_limit": limit,
	}

	data, _ := json.Marshal(payload)

	var resp GenericResponse
	err := c.doSiteRequest(http.MethodGet, site, "stat/alluser", bytes.NewReader(data), &resp)
	return &resp, err
}
