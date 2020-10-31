package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/actgardner/gogen-avro/v7/container"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

// Get the schema file from our conventional path
func LoadTestSchema() (*goavro.Codec, error) {
	schema, err := ioutil.ReadFile("schema.avsc")
	if err != nil {
		return nil, err
	}

	return goavro.NewCodec(string(schema))
}

func LoadTestFixtures() ([]json.RawMessage, error) {
	fixtureJson, err := ioutil.ReadFile("fixtures.json")
	if err != nil {
		return nil, err
	}

	fixtures := make([]json.RawMessage, 0)
	if err := json.Unmarshal([]byte(fixtureJson), &fixtures); err != nil {
		return nil, err
	}

	return fixtures, nil
}

func GGJSONToAvroBytes(fixture json.RawMessage, fixtureType container.AvroRecord) ([]byte, error) {
	if err := json.Unmarshal([]byte(fixture), &fixtureType); err != nil {
		return nil, err
	}

	var avroBytes bytes.Buffer
	if err := fixtureType.Serialize(&avroBytes); err != nil {
		return nil, err
	}
	return avroBytes.Bytes(), nil
}

func GAJSONToAvroBytes(fixture json.RawMessage, codec *goavro.Codec) ([]byte, error) {
	native, _, err := codec.NativeFromTextual([]byte(fixture))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Native: %s\n", fixture)

	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

// RoundTripExactBytes tests that:
// - the avro-encoded bytes from goavro and gogen-avro are identical
// - gogen-avro can decode avro-enocded data from goavro and the Go data is identical
// - goavro can decode JSON-encoded data from gogen-avro and the Go data is identical
//
// For schemas with maps use RoundTrip instead since maps are not encoded deterministically.
func RoundTripExactBytes(t *testing.T, recordFunc func() container.AvroRecord, deserMethod func(io.Reader) (interface{}, error)) {
	codec, err := LoadTestSchema()
	assert.NoError(t, err)

	fixtures, err := LoadTestFixtures()
	assert.NoError(t, err)

	for _, f := range fixtures {
		record := recordFunc()
		ggBytes, err := GGJSONToAvroBytes(f, record)
		assert.NoError(t, err)

		gaBytes, err := GAJSONToAvroBytes(f, codec)
		assert.NoError(t, err)

		// Confirm gogen-avro and goavro produce the same binary serialization
		assert.Equal(t, ggBytes, gaBytes)

		// Confirm that gogen-avro can deserialize the data from goavro
		deserRecord, err := deserMethod(bytes.NewBuffer(gaBytes))
		assert.NoError(t, err)
		assert.Equal(t, record, deserRecord)

		// Confirm that goavro can deserialize the JSON representation serialized by gogen-avro
		ggJson, err := json.Marshal(deserRecord)
		assert.NoError(t, err)
		gaFromJSONBytes, err := GAJSONToAvroBytes(json.RawMessage(ggJson), codec)
		assert.NoError(t, err)
		assert.Equal(t, gaBytes, gaFromJSONBytes)
	}
}

// RoundTrip tests that:
// - gogen-avro can decode avro-enocded data from goavro and the Go data is identical
// - goavro can decode JSON-encoded data from gogen-avro and the Go data is identical
//
func RoundTrip(t *testing.T, recordFunc func() container.AvroRecord, deserMethod func(io.Reader) (interface{}, error)) {
	codec, err := LoadTestSchema()
	assert.NoError(t, err)

	fixtures, err := LoadTestFixtures()
	assert.NoError(t, err)

	for _, f := range fixtures {
		record := recordFunc()
		_, err := GGJSONToAvroBytes(f, record)
		assert.NoError(t, err)

		gaBytes, err := GAJSONToAvroBytes(f, codec)
		assert.NoError(t, err)

		// Confirm that gogen-avro can deserialize the data from goavro
		deserRecord, err := deserMethod(bytes.NewBuffer(gaBytes))
		assert.NoError(t, err)
		assert.Equal(t, record, deserRecord)

		// Confirm that goavro can deserialize the JSON representation serialized by gogen-avro
		ggJson, err := json.Marshal(deserRecord)
		assert.NoError(t, err)

		gaFixture, _, err := codec.NativeFromTextual([]byte(f))
		assert.NoError(t, err)

		ggJsonNative, _, err := codec.NativeFromTextual(ggJson)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, gaFixture, ggJsonNative)
	}
}

// RoundTripGoGenOnly tests that a JSON fixture can be serialized as avro bytes, then re-serialized into equivalent JSON.
// This is used for tests that can't use goavro because the definitions are spread across multiple schema files.
func RoundTripGoGenOnly(t *testing.T, recordFunc func() container.AvroRecord, deserMethod func(io.Reader) (interface{}, error)) {
	fixtures, err := LoadTestFixtures()
	assert.NoError(t, err)

	for _, f := range fixtures {
		record := recordFunc()
		ggBytes, err := GGJSONToAvroBytes(f, record)
		assert.NoError(t, err)

		deserRecord, err := deserMethod(bytes.NewBuffer(ggBytes))
		assert.NoError(t, err)

		assert.Equal(t, record, deserRecord)

		ggJson, err := json.Marshal(deserRecord)
		assert.NoError(t, err)

		var expected, actual interface{}
		assert.NoError(t, json.Unmarshal(f, &expected))
		assert.NoError(t, json.Unmarshal(ggJson, &actual))
		assert.Equal(t, expected, actual)
	}
}
