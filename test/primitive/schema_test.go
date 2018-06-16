package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

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
		if !ok {
			t.Fatalf("GOT: %#v; WANT: %#v", ok, true)
		}
		if !reflect.DeepEqual(structVal, avroVal) {
			t.Fatalf("Field %v not equal: %v != %v", fieldName, structVal, avroVal)
		}
	}
}

func TestPrimitiveFixture(t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	if err != nil {
		t.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, remaining, err := codec.NativeFromBinary(buf.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if got, want := len(remaining), 0; got != want {
			t.Fatalf("GOT: %#v; WANT: %#v", got, want)
		}
		compareFixtureGoAvro(t, datum, f)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := DeserializePrimitiveTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}

func BenchmarkSerializePrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := record.Serialize(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSerializePrimitiveGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	if err != nil {
		b.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		b.Fatal(err)
	}
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
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeserializePrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}
	err := record.Serialize(buf)
	if err != nil {
		b.Fatal(err)
	}
	recordBytes := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DeserializePrimitiveTestRecord(bytes.NewReader(recordBytes))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeserializePrimitiveGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	if err != nil {
		b.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		b.Fatal(err)
	}
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
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := codec.NativeFromBinary(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
