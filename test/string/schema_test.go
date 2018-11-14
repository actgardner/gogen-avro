package avro

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

		buffer, err := prepareBuffer(t, length, testCase.input)
		assert.Nil(t, err)

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
func prepareBuffer(t *testing.T, length int64, input StringRec) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	err := writeLong(length, &buffer)
	assert.Nil(t, err)

	_, err = buffer.Write([]byte(input.ProductName))
	assert.Nil(t, err)
	return buffer, nil
}

func checkEqual(input StringRec, output *StringRec, t *testing.T) {
	if !reflect.DeepEqual(input, *output) {
		t.Error("deserialized content not equal to input")
	}
}
