package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SiteReport map[string]interface{}

type SiteReportsResponse struct {
	Meta CommonMeta   `json:"meta"`
	Data []SiteReport `json:"data"`
}

type ReportInterval string

const (
	ReportInterval5Min   ReportInterval = "5minutes"
	ReportIntervalHourly ReportInterval = "hourly"
	ReportIntervalDaily  ReportInterval = "daily"
)

func (r ReportInterval) Valid() bool {
	switch r {
	case ReportInterval5Min, ReportIntervalHourly, ReportIntervalDaily:
		return true
	default:
		return false
	}
}

type ReportType string

const (
	ReportTypeSite ReportType = "site"
	ReportTypeUser ReportType = "user"
	ReportTypeAP   ReportType = "ap"
)

func (r ReportType) Valid() bool {
	switch r {
	case ReportTypeSite, ReportTypeUser, ReportTypeAP:
		return true
	default:
		return false
	}
}

type ReportAttribute string

const (
	ReportAttributeBytes         ReportAttribute = "bytes"
	ReportAttributeWANTXBytes    ReportAttribute = "wan-tx_bytes"
	ReportAttributeWANRXBytes    ReportAttribute = "wan-rx_bytes"
	ReportAttributeWLANBytes     ReportAttribute = "wlan_bytes"
	ReportAttributeNumberSTA     ReportAttribute = "num_sta"
	ReportAttributeLANNumberSTA  ReportAttribute = "lan-num_sta"
	ReportAttributeWLANNumberSTA ReportAttribute = "wlan-num_sta"
	ReportAttributeTime          ReportAttribute = "time"
	ReportAttributeRXBytes       ReportAttribute = "rx_bytes"
	ReportAttributeTXBytes       ReportAttribute = "tx_bytes"
)

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

func (r ReportAttribute) Valid() bool {
	switch r {
	case ReportAttributeBytes, ReportAttributeWANTXBytes, ReportAttributeWANRXBytes, ReportAttributeWLANBytes:
		fallthrough
	case ReportAttributeNumberSTA, ReportAttributeLANNumberSTA, ReportAttributeWLANNumberSTA, ReportAttributeTime:
		fallthrough
	case ReportAttributeRXBytes, ReportAttributeTXBytes:
		return true
	default:
		return false
	}
}

func (r ReportAttribute) MarshalJSON() ([]byte, error) {
	return []byte(string(r)), nil
}

func (r *ReportAttribute) UnmarshalJSON(data []byte) error {
	*r = ReportAttribute(bytes.NewBuffer(data).String())
	return nil
}

func (c *Client) SiteReport(site string, interval ReportInterval, reportType ReportType, attributesToReturn []ReportAttribute, filterMacs ...string) (*SiteReportsResponse, error) {
	if !interval.Valid() {
		return nil, fmt.Errorf("invalid interval specified: %s", interval)
	}
	if !reportType.Valid() {
		return nil, fmt.Errorf("invalid reportType specified: %s", reportType)
	}

	if len(attributesToReturn) == 0 {
		attributesToReturn = AllReportAttributes
	} else {
		for _, attr := range attributesToReturn {
			if !attr.Valid() {
				return nil, fmt.Errorf("invalid report attribute specified: %s", attr)
			}
		}
	}

	payload := map[string]interface{}{
		"attributes": attributesToReturn,
	}
	if len(filterMacs) > 0 {
		payload["macs"] = filterMacs
	}

	data, _ := json.Marshal(payload)

	var resp SiteReportsResponse
	err := c.doSiteRequest(http.MethodGet, site, fmt.Sprintf("stat/report/%s.%s", interval, reportType), bytes.NewReader(data), &resp)
	return &resp, err
}
