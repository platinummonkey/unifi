package unifi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
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

var zeroValueReflect = reflect.Value{}

// UnmarshalJSON will unmarshal JSON data and catch extra fields into the XXXUnknown
func UnmarshalJSON(data []byte, obj interface{}) error {
	objValue := reflect.ValueOf(obj).Elem()
	knownFields := make(map[string]reflect.Value, 0)
	xField := objValue.FieldByName("XXXUnknown")
	for i := 0; i != objValue.NumField(); i++ {
		if objValue.Type().Field(i).Name != "XXXUnknown" {
			jsonTag := objValue.Type().Field(i).Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			if jsonTag == "" {
				jsonTag = objValue.Type().Field(i).Name
			}
			jsonName := strings.Split(jsonTag, ",")[0]
			knownFields[jsonName] = objValue.Field(i)
		}
	}

	if xField == zeroValueReflect || xField.Kind() != reflect.Map {
		// no extra field, just unmarshal directly.
		log.Printf("no XXXUnknown field: %T %v %v\n", xField, xField, xField.Kind())
		return json.Unmarshal(data, obj)
	}

	xFieldContainer := make(map[string]interface{}, 0)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(&xFieldContainer)
	if err != nil {
		return err
	}

	for knownField, valueContainer := range knownFields {
		if val, ok := xFieldContainer[knownField]; ok {
			switch valueContainer.Kind() {
			case reflect.Int:
				if jsonVal, ok := val.(json.Number); ok {
					intVal, jErr := jsonVal.Int64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(int(intVal)))
					}
				}
			case reflect.Int8:
				if jsonVal, ok := val.(json.Number); ok {
					intVal, jErr := jsonVal.Int64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(int8(intVal)))
					}
				}
			case reflect.Int16:
				if jsonVal, ok := val.(json.Number); ok {
					intVal, jErr := jsonVal.Int64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(int16(intVal)))
					}
				}
			case reflect.Int32:
				if jsonVal, ok := val.(json.Number); ok {
					intVal, jErr := jsonVal.Int64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(int32(intVal)))
					}
				}
			case reflect.Int64:
				if jsonVal, ok := val.(json.Number); ok {
					intVal, jErr := jsonVal.Int64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(int64(intVal)))
					}
				}
			case reflect.Uint:
				if jsonVal, ok := val.(json.Number); ok {
					uintVal, jErr := strconv.ParseUint(jsonVal.String(), 10, 64)
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(uint(uintVal)))
					}
				}
			case reflect.Uint8:
				if jsonVal, ok := val.(json.Number); ok {
					uintVal, jErr := strconv.ParseUint(jsonVal.String(), 10, 64)
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(uint8(uintVal)))
					}
				}
			case reflect.Uint16:
				if jsonVal, ok := val.(json.Number); ok {
					uintVal, jErr := strconv.ParseUint(jsonVal.String(), 10, 64)
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(uint16(uintVal)))
					}
				}
			case reflect.Uint32:
				if jsonVal, ok := val.(json.Number); ok {
					uintVal, jErr := strconv.ParseUint(jsonVal.String(), 10, 64)
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(uint32(uintVal)))
					}
				}
			case reflect.Uint64:
				if jsonVal, ok := val.(json.Number); ok {
					uintVal, jErr := strconv.ParseUint(jsonVal.String(), 10, 64)
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(uint64(uintVal)))
					}
				}
			case reflect.Float32:
				if jsonVal, ok := val.(json.Number); ok {
					floatVal, jErr := jsonVal.Float64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(float32(floatVal)))
					}
				}
			case reflect.Float64:
				if jsonVal, ok := val.(json.Number); ok {
					floatVal, jErr := jsonVal.Float64()
					if jErr == nil {
						valueContainer.Set(reflect.ValueOf(float64(floatVal)))
					}
				}
			default:
				if _, ok := val.(json.Number); !ok {
					// maybe a complex number otherwise? ignore for now, not supported on purpose
					valueContainer.Set(reflect.ValueOf(val))
				}
			}
		}
		delete(xFieldContainer, knownField)
	}

	// turn back json.Number to float64
	fixedFields := make(map[string]interface{}, 0)
	for key, val := range xFieldContainer {
		if jsonVal, ok := val.(json.Number); ok {
			intVal, jErr := jsonVal.Int64()
			if jErr == nil {
				fixedFields[key] = intVal
			} else {
				floatVal, jErr := jsonVal.Float64()
				if jErr == nil {
					fixedFields[key] = floatVal
				}
			}
		}
	}
	for key, value := range fixedFields {
		xFieldContainer[key] = value
	}
	xField.Set(reflect.ValueOf(xFieldContainer))

	return nil
}

// MarshalJSON is a wrapper around json.Marshal to compile XXXUnknown fields into the output response.
func MarshalJSON(obj interface{}) ([]byte, error) {
	objValue := reflect.ValueOf(obj).Elem()
	xField := objValue.FieldByName("XXXUnknown")
	if xField == zeroValueReflect || xField.Kind() != reflect.Map {
		// no extra field, just unmarshal directly.
		log.Printf("no X field: %T %v %v\n", xField, xField, xField.Kind())
		return json.Marshal(obj)
	}

	output := make(map[string]interface{}, 0)
	knownFields := make(map[string]struct{}, 0)

	for i := 0; i != objValue.NumField(); i++ {
		if objValue.Type().Field(i).Name != "XXXUnknown" {
			jsonTag := objValue.Type().Field(i).Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			if jsonTag == "" {
				jsonTag = objValue.Type().Field(i).Name
			}
			jsonName := strings.Split(jsonTag, ",")[0]
			output[jsonName] = objValue.Field(i).Interface()
			knownFields[jsonName] = struct{}{}
		} else {
			mapIter := xField.MapRange()
			for mapIter.Next() {
				if mapIter.Key().Kind() == reflect.String {
					output[mapIter.Key().String()] = mapIter.Value().Interface()
				}
			}
		}
	}

	return json.Marshal(output)
}

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
