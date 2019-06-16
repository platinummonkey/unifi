package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SiteReport is the site report data
/*type SiteReport struct {
	// common
	OID string `json:"oid"`
	// site report
	Site string `json:"site"`
	// ap report
	AccessPoint string `json:"ap"`
	// user
	???
	// speedtest
	SpeedTestID string `json:"_id"`
    Origin string `json:"o"`
}*/
type SiteReport map[string]interface{}

// SiteReportsResponse contains the site report data
type SiteReportsResponse struct {
	Meta CommonMeta   `json:"meta"`
	Data []SiteReport `json:"data"`
}

// ReportInterval is a defined report interval
// there are only a few known intervals supported
type ReportInterval string

// The supported Report Intervals
const (
	ReportInterval5Min    ReportInterval = "5minutes"
	ReportIntervalHourly  ReportInterval = "hourly"
	ReportIntervalDaily   ReportInterval = "daily"
	ReportIntervalArchive ReportInterval = "archive"
)

// IsValid returns true if it's a valid report interval.
// there are only a few valid types
func (r ReportInterval) IsValid() bool {
	switch r {
	case ReportInterval5Min, ReportIntervalHourly, ReportIntervalDaily, ReportIntervalArchive:
		return true
	default:
		return false
	}
}

// ReportType defines the report type
type ReportType string

// The supported Report Types
const (
	ReportTypeSite      ReportType = "site"
	ReportTypeUser      ReportType = "user"
	ReportTypeAP        ReportType = "ap"
	ReportTypeSpeedTest ReportType = "speedtest"
)

// IsValid returns true if it's a valid report type.
// there are only a few valid types
func (r ReportType) IsValid() bool {
	switch r {
	case ReportTypeSite, ReportTypeUser, ReportTypeAP, ReportTypeSpeedTest:
		return true
	default:
		return false
	}
}

// ReportAttribute defines the report attribute
type ReportAttribute string

// Available report attributes
const (
	ReportAttributeBytes             ReportAttribute = "bytes"
	ReportAttributeWANTXBytes        ReportAttribute = "wan-tx_bytes"
	ReportAttributeWANRXBytes        ReportAttribute = "wan-rx_bytes"
	ReportAttributeWLANBytes         ReportAttribute = "wlan_bytes"
	ReportAttributeNumberSTA         ReportAttribute = "num_sta"
	ReportAttributeLANNumberSTA      ReportAttribute = "lan-num_sta"
	ReportAttributeWLANNumberSTA     ReportAttribute = "wlan-num_sta"
	ReportAttributeTime              ReportAttribute = "time"
	ReportAttributeRXBytes           ReportAttribute = "rx_bytes"
	ReportAttributeTXBytes           ReportAttribute = "tx_bytes"
	ReportAttributeSpeedTestDownload ReportAttribute = "xput_download"
	ReportAttributeSpeedTestUpload   ReportAttribute = "xput_upload"
	ReportAttributeSpeedTestLatency  ReportAttribute = "latency"
)

// AllReportAttributes contains all the normal report attributes
// the only known difference is the speed test report which should use SpeedTestReportAttributes
var AllReportAttributes = []ReportAttribute{
	ReportAttributeBytes,
	ReportAttributeWANTXBytes,
	ReportAttributeWANRXBytes,
	ReportAttributeWLANBytes,
	ReportAttributeNumberSTA,
	ReportAttributeLANNumberSTA,
	ReportAttributeWLANNumberSTA,
	ReportAttributeTime,
	ReportAttributeRXBytes,
	ReportAttributeTXBytes,
}

// SpeedTestReportAttributes contains all the report attributes for the speed-test report
var SpeedTestReportAttributes = []ReportAttribute{
	ReportAttributeSpeedTestDownload,
	ReportAttributeSpeedTestUpload,
	ReportAttributeSpeedTestLatency,
	ReportAttributeTime,
}

// IsValid returns true if it's a valid report attribute.
// there are only a few valid types
func (r ReportAttribute) IsValid() bool {
	switch r {
	case ReportAttributeBytes, ReportAttributeWANTXBytes, ReportAttributeWANRXBytes, ReportAttributeWLANBytes:
		fallthrough
	case ReportAttributeNumberSTA, ReportAttributeLANNumberSTA, ReportAttributeWLANNumberSTA, ReportAttributeTime:
		fallthrough
	case ReportAttributeRXBytes, ReportAttributeTXBytes, ReportAttributeSpeedTestDownload:
		fallthrough
	case ReportAttributeSpeedTestUpload, ReportAttributeSpeedTestLatency:
		return true
	default:
		return false
	}
}

// MarshalJSON implements json.Marshaler
func (r ReportAttribute) MarshalJSON() ([]byte, error) {
	return []byte(string(r)), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (r *ReportAttribute) UnmarshalJSON(data []byte) error {
	*r = ReportAttribute(bytes.NewBuffer(data).String())
	return nil
}

// SiteReport returns the site stats method for the given report interval and type of report
// site - the site interested in stats
// startTime - start time of the report, set to 0 and endTime to 0 for default behavior
// endTime - end time of the report, set to 0 and startTime to 0 for default behavior
// interval - the report interval requested
// reportType - the report type requested
// attributes - attributes to return, see AllReportAttributes for default behavior
// filterMacs - optional list of macs to filter stats.
func (c *Client) SiteReport(site string, startTime time.Time, endTime time.Time, interval ReportInterval, reportType ReportType, attributes []ReportAttribute, filterMacs ...string) (*SiteReportsResponse, error) {
	if startTime.IsZero() && endTime.IsZero() {
		endTime := time.Now().UTC()
		switch interval {
		case ReportInterval5Min:
			// set default to last 1h
			startTime = endTime.Add(-1 * time.Hour)
		case ReportIntervalHourly:
			// set default to last 1 day
			startTime = endTime.Add(-24 * time.Hour)
		case ReportIntervalDaily:
			// set default to last 7 days
			startTime = endTime.Add(7 * 24 * time.Hour)
		}
	}

	if !startTime.Before(endTime) || startTime == endTime {
		return nil, fmt.Errorf("invalid end time, must occur after start time")
	}

	if !reportType.IsValid() {
		return nil, fmt.Errorf("invalid reportType specified: %s", reportType)
	}
	// only archive is supported for speedtest, so override.
	if reportType == ReportTypeSpeedTest {
		interval = ReportIntervalArchive
	}

	if !interval.IsValid() {
		return nil, fmt.Errorf("invalid interval specified: %s", interval)
	}

	if len(attributes) == 0 {
		attributes = AllReportAttributes
		if reportType == ReportTypeSpeedTest {
			attributes = SpeedTestReportAttributes
		}
	} else {
		for _, attr := range attributes {
			if !attr.IsValid() {
				return nil, fmt.Errorf("invalid report attribute specified: %s", attr)
			}
		}
	}

	payload := map[string]interface{}{
		"attributes": attributes,
		"start":      startTime.UTC().Unix() * 1000,
		"end":        startTime.UTC().Unix() * 1000,
	}
	if len(filterMacs) > 0 {
		payload["macs"] = filterMacs
	}

	data, _ := json.Marshal(payload)

	var resp SiteReportsResponse
	err := c.doSiteRequest(http.MethodGet, site, fmt.Sprintf("stat/report/%s.%s", interval, reportType), bytes.NewReader(data), &resp)
	return &resp, err
}
