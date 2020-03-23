package avro

import (
	"bytes"
	"encoding/json"
	"github.com/actgardner/gogen-avro/singleobject"
	"io/ioutil"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
	{
		"EnumField": 0
	},
	{
		"EnumField": 2
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
		writer := singleobject.NewWriter(&buf, EnumTestRecordAvroCRC64Fingerprint)
		err = f.Serialize(writer)
		assert.Nil(t, err)
		b := singleobject.NewReader(&buf).Bytes()
		datum, remaining, err := codec.NativeFromBinary(b)
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
		writer := singleobject.NewWriter(&buf, EnumTestRecordAvroCRC64Fingerprint)
		err = f.Serialize(writer)
		assert.Nil(t, err)

		datum, err := DeserializeEnumTestRecord(singleobject.NewReader(&buf))
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}

func TestInvalidStringConversion(t *testing.T) {
	enumified, err := NewTestEnumTypeValue("bogus")
	assert.Error(t, err)
	assert.Equal(t, TestEnumType(-1), enumified)
}
