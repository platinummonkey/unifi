package unifi

import (
	"fmt"
	"strings"
)

const (
	ContentTypeHeader = "application/json"
	Version           = "0.0.1"
	UserAgentHeader   = "unifi/" + Version
)

// common errors
const (
	APINoPermissionError = "api.err.NoPermission"
	APIInvalidError      = "api.err.Invalid"
)

// InvalidResponseBody indicates and error with the body of the response
var InvalidResponseBody = fmt.Errorf("invalid response body")

// JSONDecodeError indicates an unexpected unmarshal problem from the API, check this is a valid endpoint.
var JSONDecodeError = fmt.Errorf("unable to unmarshal json response")

type ResponseCode string

func (r ResponseCode) MarshalJSON() ([]byte, error) {
	return []byte(string(r)), nil
}

func (r *ResponseCode) UnmarshalJSON(data []byte) error {
	// HACK: we must trim `"` because sometimes the controller passes back `"rc": "\"<status>\""`
	*r = ResponseCode(strings.TrimRight(strings.TrimLeft(string(data), "\""), "\""))
	return nil
}

func (r *ResponseCode) Equal(o ResponseCode) bool {
	return strings.EqualFold(string(*r), string(o))
}

const (
	ResponseCodeOK    ResponseCode = "ok"
	ResponseCodeError ResponseCode = "error"
)

type CommonMeta struct {
	ResponseCode        ResponseCode `json:"rc"`
	ResponseCodeMessage string       `json:"msg"`
	Count               int          `json:"count"`

	XXXUnknown map[string]interface{} `json:"-"`
}

func (m *CommonMeta) GetResponseCode() ResponseCode {
	return m.ResponseCode
}

func (m *CommonMeta) GetResponseMessage() string {
	return m.ResponseCodeMessage
}

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

type CallableCommand interface {
	Manager() string                         // must return the manager
	Command() string                         // must return the command
	Payload() (map[string]interface{}, bool) // must return any additional payload, false if nothing to include
}

type GenericResponse struct {
	Meta CommonMeta             `json:"meta"`
	Data map[string]interface{} `json:"data"`
}
