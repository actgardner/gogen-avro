package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/actgardner/gogen-avro/v8/container"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

type RecordFactory func() container.AvroRecord
type DeserMethod func(io.Reader) (container.AvroRecord, error)
type EvolutionDeserMethod func(io.Reader, string) (container.AvroRecord, error)

type EvolutionFixture struct {
	Data     json.RawMessage
	Expected json.RawMessage
	Err      *string
}

// Get the schema file from our conventional path
func LoadTestSchema() ([]byte, error) {
	return ioutil.ReadFile("schema.avsc")
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

func LoadEvolutionFixtures() ([]EvolutionFixture, error) {
	fixtureJson, err := ioutil.ReadFile("evolution.json")
	if err != nil {
		return nil, err
	}

	fixtures := make([]EvolutionFixture, 0)
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
func RoundTripExactBytes(t *testing.T, recordFunc RecordFactory, deserMethod DeserMethod) {
	schema, err := LoadTestSchema()
	assert.NoError(t, err)

	codec, err := goavro.NewCodec(string(schema))
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
func RoundTrip(t *testing.T, recordFunc RecordFactory, deserMethod DeserMethod) {
	schema, err := LoadTestSchema()
	assert.NoError(t, err)

	codec, err := goavro.NewCodec(string(schema))
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
func RoundTripGoGenOnly(t *testing.T, recordFunc RecordFactory, deserMethod DeserMethod) {
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

func RoundTripEvolution(t *testing.T, oldRecordFunc, newRecordFunc RecordFactory, newDeserMethod EvolutionDeserMethod) {
	oldSchema, err := LoadTestSchema()
	assert.NoError(t, err)

	fixtures, err := LoadEvolutionFixtures()
	assert.NoError(t, err)

	for _, f := range fixtures {
		// Serialize the fixture into Avro bytes using the old schema
		oldRecord := oldRecordFunc()
		oldBytes, err := GGJSONToAvroBytes(f.Data, oldRecord)
		assert.NoError(t, err)

		// Deserialize the Avro data with the new schema and compare to the expected JSON deserialization
		newRecord, err := newDeserMethod(bytes.NewBuffer(oldBytes), string(oldSchema))
		assert.NoError(t, err)

		expectedRecord := newRecordFunc()
		err = json.Unmarshal([]byte(f.Expected), &expectedRecord)
		assert.NoError(t, err)

		assert.Equal(t, expectedRecord, newRecord)

		// Deserialize the JSON fixture with the new schema
		jsonRecord := newRecordFunc()
		err = json.Unmarshal([]byte(f.Data), &jsonRecord)
		assert.NoError(t, err)
		assert.Equal(t, expectedRecord, jsonRecord)
	}
}
