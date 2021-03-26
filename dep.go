package unifi

import (
	"fmt"
	"strings"
)

// Common headers
const (
	ContentTypeHeader = "application/json"
	Version           = "0.0.1"
	UserAgentHeader   = "unifi/" + Version
)

// Common errors
const (
	APINoPermissionError = "api.err.NoPermission"
	APIInvalidError      = "api.err.Invalid"
)

// ErrInvalidResponseBody indicates and error with the body of the response
var ErrInvalidResponseBody = fmt.Errorf("invalid response body")

// ErrJSONDecode indicates an unexpected unmarshal problem from the API, check this is a valid endpoint.
var ErrJSONDecode = fmt.Errorf("unable to unmarshal json response")

// ResponseCode is the api response code, typically just `ok` or `err`
type ResponseCode string

// MarshalJSON implements json.Marshaler
func (r ResponseCode) MarshalJSON() ([]byte, error) {
	return []byte(string(r)), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (r *ResponseCode) UnmarshalJSON(data []byte) error {
	// HACK: we must trim `"` because sometimes the controller passes back `"rc": "\"<status>\""`
	*r = ResponseCode(strings.TrimRight(strings.TrimLeft(string(data), "\""), "\""))
	return nil
}

// Equal compares response codes
func (r *ResponseCode) Equal(o ResponseCode) bool {
	return strings.EqualFold(
		strings.TrimSpace(strings.Trim(strings.ToLower(string(*r)), "\"")),
		strings.TrimSpace(strings.Trim(strings.ToLower(string(o)), "\"")),
	)
}

// Known response codes from the api
const (
	ResponseCodeOK    ResponseCode = "ok"
	ResponseCodeError ResponseCode = "error"
)

// CommonMeta is the most common meta response from the api
type CommonMeta struct {
	ResponseCode        ResponseCode `json:"rc"`
	ResponseCodeMessage string       `json:"msg"`
	Count               int          `json:"count"`

	XXXUnknown map[string]interface{} `json:"-"`
}

// GetResponseCode returns the response code
func (m *CommonMeta) GetResponseCode() ResponseCode {
	return m.ResponseCode
}

// GetResponseMessage return the response message
func (m *CommonMeta) GetResponseMessage() string {
	return m.ResponseCodeMessage
}

// GeoCodeData contains the geo-code data that is common
type GeoCodeData struct {
	AreaCode      int     `json:"area_code"`
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	CountryCode   string  `json:"country_code"`
	CountryCode3  string  `json:"country_code3"`
	CountryName   string  `json:"country_name"`
	DMACode       int     `json:"dma_code"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	PostalCode    string  `json:"postal_code"` // varies based on country
	Region        string  `json:"region"`
}

// GenericResponse is the most generic response
// this is used in the short term to provide quick functionality while they are more defined in experimentation/docs.
type GenericResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data []map[string]interface{} `json:"data"`
}
