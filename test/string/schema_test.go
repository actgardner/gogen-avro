package avro

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"testing"
)

var testCases = []struct {
	input         StringRec
	length        int64
	expectFailure bool
}{
	{StringRec{"Test string"}, 0, false},
	{StringRec{"Test string"}, math.MaxInt64, true},
	{StringRec{"Test string"}, -1, true},
}

// Test string deserializer behaviour to check that corrupted string will not generate a panic.
func TestCorruptString(t *testing.T) {
	for _, testCase := range testCases {
		var length int64 = testCase.length

		// If length is 0, calculate actual string length, otherwise take test value to "corrupt" it.
		if testCase.length == 0 {
			length = int64(len(testCase.input.ProductName))
		}

		buffer, err := prepareBuffer(length, testCase.input)
		if err != nil {
			t.Error("cannot prepare buffer:", err)
		}

		output, err := DeserializeStringRec(&buffer)
		if err != nil {
			if !testCase.expectFailure {
				t.Error("deserialize failed:", err)
			}
			continue
		}

		if testCase.expectFailure {
			t.Error("expecting error on deserialize")
			continue
		}

		checkEqual(testCase.input, output, t)
	}
}

// Manually prepare simple Avro string buffer to generate possibly corrupted string.
func prepareBuffer(length int64, input StringRec) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	err := writeLong(length, &buffer)
	if err != nil {
		return buffer, errors.New("cannot generate string length")
	}

	_, err = buffer.Write([]byte(input.ProductName))
	if err != nil {
		return buffer, errors.New("cannot append string")
	}
	return buffer, nil
}

func checkEqual(input StringRec, output *StringRec, t *testing.T) {
	if !reflect.DeepEqual(input, *output) {
		t.Error("deserialized content not equal to input")
	}
}
