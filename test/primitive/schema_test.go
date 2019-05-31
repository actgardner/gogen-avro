package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/schema"
	"github.com/actgardner/gogen-avro/vm"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
{"IntField": 1, "LongField": 2, "FloatField": 3.4, "DoubleField": 5.6, "StringField": "789", "BoolField": true, "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"},
{"IntField": 2147483647, "LongField": 9223372036854775807, "FloatField": 3.402823e+38, "DoubleField": 1.7976931348623157e+308, "StringField": "", "BoolField": false, "BytesField": ""},
{"IntField": -2147483647, "LongField": -9223372036854775807, "FloatField": 3.402823e-38, "DoubleField": 2.2250738585072014e-308, "StringField": "", "BoolField": true, "BytesField": ""}
]
`

func compareFixtureGoAvro(t *testing.T, actual interface{}, expected interface{}) {
	record := actual.(map[string]interface{})
	value := reflect.ValueOf(expected)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		structVal := value.Field(i).Interface()
		avroVal, ok := record[fieldName]
		assert.Equal(t, true, ok)
		assert.Equal(t, structVal, avroVal)
	}
}

func TestPrimitiveFixture(t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	schemaJson, err := ioutil.ReadFile("primitives.avsc")
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
		compareFixtureGoAvro(t, datum, f)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		target, err := DeserializePrimitiveTestRecord(&buf)
		assert.Nil(t, err)

		assert.Equal(t, target, &f)
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

	readerNs := schema.NewNamespace(false)
	readerType, err := readerNs.TypeForSchema(schemaJson)
	assert.Nil(b, err)

	err = readerType.ResolveReferences(readerNs)
	assert.Nil(b, err)

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
