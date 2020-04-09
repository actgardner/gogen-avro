package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
	{
		"EnumField": "TestSymbol1"
	},
	{
		"EnumField": "testSymbol2"
	}
]
`

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
	}
}

func TestInvalidStringConversion(t *testing.T) {
	enumified, err := NewTestEnumTypeValue("bogus")
	assert.Error(t, err)
	assert.Equal(t, TestEnumType(-1), enumified)
}
