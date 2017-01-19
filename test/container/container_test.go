package avro

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	"reflect"
	"testing"
)

/* Round-trip some primitive values through our container file writer and goavro to verify */
const fixtureJson = `
[
{"IntField": 1, "LongField": 2, "FloatField": 3.4, "DoubleField": 5.6, "StringField": "789", "BoolField": true, "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"},
{"IntField": 2147483647, "LongField": 9223372036854775807, "FloatField": 3.402823e+38, "DoubleField": 1.7976931348623157e+308, "StringField": "", "BoolField": false, "BytesField": ""},
{"IntField": -2147483647, "LongField": -9223372036854775807, "FloatField": 3.402823e-38, "DoubleField": 2.2250738585072014e-308, "StringField": "", "BoolField": true, "BytesField": ""}
]
`

func TestNullEncoding(t *testing.T) {
	roundTripWithCodec(Null, t)
}

func TestSnappyEncoding(t *testing.T) {
	roundTripWithCodec(Deflate, t)
}

func TestDeflateEncoding(t *testing.T) {
	roundTripWithCodec(Snappy, t)
}

func roundTripWithCodec(codec Codec, t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	// Write the container file contents to the buffer
	var containerWriter *PrimitiveTestRecordContainerWriter
	containerWriter, err = NewPrimitiveTestRecordContainerWriter(&buf, codec, 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fixtures {
		// Write the record to the container file
		err = containerWriter.WriteRecord(f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}

	reader, err := goavro.NewReader(goavro.FromReader(&buf))
	if err != nil {
		t.Fatal(err)
	}

	var i int
	for reader.Scan() {
		datum, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}
		record := datum.(*goavro.Record)
		value := reflect.ValueOf(fixtures[i])
		for j := 0; j < value.NumField(); j++ {
			fieldName := value.Type().Field(j).Name
			structVal := value.Field(j).Interface()
			avroVal, err := record.Get(fieldName)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(structVal, avroVal) {
				t.Fatalf("Field %v not equal: %v != %v", fieldName, structVal, avroVal)
			}
		}
		i = i + 1
	}
}
