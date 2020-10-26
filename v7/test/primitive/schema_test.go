package avro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/container"
	"github.com/actgardner/gogen-avro/v7/parser"
	"github.com/actgardner/gogen-avro/v7/resolver"
	"github.com/actgardner/gogen-avro/v7/vm"

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

// Deserialize an JSON-encoded Avro payload using goavro and return the Avro-encoded bytes and the native representation
func GADeserializeFixture(fixture json.RawMessage, codec *goavro.Codec) ([]byte, error) {
	native, _, err := codec.NativeFromTextual([]byte(fixture))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Native: %v\n", native)

	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func TestRoundTrips(t *testing.T) {
	codec, fixtures, err := LoadTestData()
	assert.NoError(t, err)

	for _, f := range fixtures {
		var record PrimitiveTestRecord
		ggBytes, err := GGDeserializeFixture(f, &record)
		assert.NoError(t, err)

		gaBytes, err := GADeserializeFixture(f, codec)
		assert.NoError(t, err)

		// Confirm gogen-avro and goavro serialize the data exactly the same
		assert.Equal(t, ggBytes, gaBytes)

		// Confirm that gogen-avro can deserialize the data from goavro
		deserRecord, err := DeserializePrimitiveTestRecord(bytes.NewBuffer(gaBytes))
		assert.NoError(t, err)
		assert.Equal(t, &record, deserRecord)
	}
}

func BenchmarkSerializePrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := record.Serialize(buf)
		assert.Nil(b, err)
	}
}

func BenchmarkSerializePrimitiveGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"IntField":    int32(1),
		"LongField":   int64(2),
		"FloatField":  float32(3.4),
		"DoubleField": float64(5.6),
		"StringField": "789",
		"BoolField":   true,
		"BytesField":  []byte{1, 2, 3, 4},
	}
	buf := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := codec.BinaryFromNative(buf, someRecord)
		assert.Nil(b, err)
	}
}

func BenchmarkDeserializePrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}
	err := record.Serialize(buf)
	assert.Nil(b, err)

	recordBytes := buf.Bytes()

	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	assert.Nil(b, err)

	readerNs := parser.NewNamespace(false)
	readerType, err := readerNs.TypeForSchema(schemaJson)
	assert.Nil(b, err)

	for _, def := range readerNs.Roots {
		err = resolver.ResolveDefinition(def, readerNs.Definitions)
		assert.Nil(b, err)
	}

	deser, err := compiler.Compile(readerType, readerType)
	assert.Nil(b, err)

	var target PrimitiveTestRecord

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := vm.Eval(bytes.NewReader(recordBytes), deser, &target)
		assert.Nil(b, err)
	}
}

func BenchmarkDeserializePrimitiveGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"IntField":    int32(1),
		"LongField":   int64(2),
		"FloatField":  float32(3.4),
		"DoubleField": float64(5.6),
		"StringField": "789",
		"BoolField":   true,
		"BytesField":  []byte{1, 2, 3, 4},
	}

	buf, err := codec.BinaryFromNative(nil, someRecord)
	assert.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := codec.NativeFromBinary(buf)
		assert.Nil(b, err)
	}
}
