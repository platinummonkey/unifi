package unifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	asrt := assert.New(t)

	type testInputJSON struct {
		A int `json:"a"`
		B string `json:"b"`
		C bool `json:"c"`
		XXXUnknown map[string]interface{} `json:"-"`
	}
	var testInput testInputJSON
	input := []byte(`{"a": 1, "b": "foo", "c": true, "d": 2, "e": "bar", "f": false, "g": 3.14}`)
	err := UnmarshalJSON(input, &testInput)
	asrt.NoError(err)
	asrt.Equal(1, testInput.A)
	asrt.Equal("foo", testInput.B)
	asrt.Equal(true, testInput.C)
	asrt.EqualValues(map[string]interface{}{"d": int64(2), "e": "bar", "f": false, "g": float64(3.14)}, testInput.XXXUnknown)
}

func TestMarshalJSON(t *testing.T) {
	asrt := assert.New(t)

	type testOutputJSON struct {
		A int `json:"a"`
		B string `json:"b"`
		C bool `json:"c"`
		XXXUnknown map[string]interface{} `json:"-"`
	}
	testOutput := testOutputJSON{
		A: 1,
		B: "foo",
		C: true,
		XXXUnknown: map[string]interface{}{
			"d": int64(2),
			"e": "bar",
			"f": false,
			"g": float64(3.14),
		},
	}
	data, err := MarshalJSON(&testOutput)
	asrt.NoError(err)
	asrt.NotEmpty(data)
	var testInput testOutputJSON
	err = UnmarshalJSON(data, &testInput)
	asrt.NoError(err)

	asrt.Equal(testOutput.A, testInput.A)
	asrt.Equal(testOutput.B, testInput.B)
	asrt.Equal(testOutput.C, testInput.C)
	asrt.EqualValues(testOutput.XXXUnknown, testInput.XXXUnknown)
}