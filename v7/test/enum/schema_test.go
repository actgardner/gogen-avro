package avro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
	{
		"EnumField": "TestSymbol1"
	},
	{
		"EnumField": "testSymbol3"
	}
]
`

func TestMarshalUnmarshal(t *testing.T) {
	expected := EnumTestRecord{EnumField: TestEnumTypeTestSymbol3}
	bytes, _ := json.Marshal(&expected)
	fmt.Printf("JSON: %s\n", bytes)
	var result EnumTestRecord
	json.Unmarshal(bytes, &result)
	assert.Equal(t, expected, result)
}

func TestEnumFixture(t *testing.T) {
	fixtures := make([]EnumTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	schemaJson, err := ioutil.ReadFile("enum.avsc")
	assert.Nil(t, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, remaining, err := codec.NativeFromBinary(buf.Bytes())
		assert.Nil(t, err)
		assert.Equal(t, 0, len(remaining))

		record := datum.(map[string]interface{})
		recordVal, ok := record["EnumField"]
		assert.Equal(t, true, ok)
		assert.Equal(t, recordVal, f.EnumField.String())

		enumified, err := NewTestEnumTypeValue(f.EnumField.String())
		assert.Nil(t, err)
		assert.Equal(t, f.EnumField, enumified)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]EnumTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeEnumTestRecord(&buf)
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
		fmt.Printf("Datum: %v %v\n", datum, f)
	}
}

func TestInvalidStringConversion(t *testing.T) {
	enumified, err := NewTestEnumTypeValue("bogus")
	assert.Error(t, err)
	assert.Equal(t, TestEnumType(-1), enumified)
}
