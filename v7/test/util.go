package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/actgardner/gogen-avro/v7/container"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

// Get the schema file and test fixture data from our conventional paths
func LoadTestData() (*goavro.Codec, []json.RawMessage, error) {
	schema, err := ioutil.ReadFile("schema.avsc")
	if err != nil {
		return nil, nil, err
	}

	codec, err := goavro.NewCodec(string(schema))
	if err != nil {
		return nil, nil, err
	}

	fixtureJson, err := ioutil.ReadFile("fixtures.json")
	if err != nil {
		return nil, nil, err
	}

	fixtures := make([]json.RawMessage, 0)
	if err := json.Unmarshal([]byte(fixtureJson), &fixtures); err != nil {
		return nil, nil, err
	}

	return codec, fixtures, nil
}

// Deserialize an JSON-encoded Avro payload using gogen-avro and return the Avro-encoded bytes and native representation
func GGDeserializeFixture(fixture json.RawMessage, fixtureType container.AvroRecord) ([]byte, error) {
	if err := json.Unmarshal([]byte(fixture), &fixtureType); err != nil {
		return nil, err
	}

	var avroBytes bytes.Buffer
	if err := fixtureType.Serialize(&avroBytes); err != nil {
		return nil, err
	}
	return avroBytes.Bytes(), nil
}

// Deserialize an JSON-encoded Avro payload using goavro and return the Avro-encoded bytes
func GADeserializeFixture(fixture json.RawMessage, codec *goavro.Codec) ([]byte, error) {
	native, _, err := codec.NativeFromTextual([]byte(fixture))
	if err != nil {
		return nil, err
	}

	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func RoundTrip(t *testing.T, record container.AvroRecord, deserMethod func(io.Reader) (interface{}, error)) {
	codec, fixtures, err := LoadTestData()
	assert.NoError(t, err)

	for _, f := range fixtures {
		ggBytes, err := GGDeserializeFixture(f, record)
		assert.NoError(t, err)

		gaBytes, err := GADeserializeFixture(f, codec)
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
		gaFromJSONBytes, err := GADeserializeFixture(json.RawMessage(ggJson), codec)
		assert.NoError(t, err)
		assert.Equal(t, gaBytes, gaFromJSONBytes)
	}
}
