package avro

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	avrotest "github.com/alanctgardner/gogen-avro/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
{"IntField": 1, "LongField": 2, "FloatField": 3.4, "DoubleField": 5.6, "StringField": "789", "BoolField": true, "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"},
{"IntField": 2147483647, "LongField": 9223372036854775807, "FloatField": 3.402823e+38, "DoubleField": 1.7976931348623157e+308, "StringField": "", "BoolField": false, "BytesField": ""},
{"IntField": -2147483647, "LongField": -9223372036854775807, "FloatField": 3.402823e-38, "DoubleField": 2.2250738585072014e-308, "StringField": "", "BoolField": true, "BytesField": ""}
]
`

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

	for _, f := range fixtures {
		avrotest.RoundTripGoAvroTest(f, codec, t)
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
	for i := 0; i < b.N; i++ {
		record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}
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
	someRecord, err := goavro.NewRecord(goavro.RecordSchema(string(schemaJson)))
	if err != nil {
		b.Fatal(err)
	}
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		someRecord.Set("IntField", int32(1))
		someRecord.Set("LongField", int64(2))
		someRecord.Set("FloatField", float32(3.4))
		someRecord.Set("DoubleField", float64(5.6))
		someRecord.Set("StringField", "789")
		someRecord.Set("BoolField", true)
		someRecord.Set("BytesField", []byte{1, 2, 3, 4})

		err := codec.Encode(buf, someRecord)
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
	for i := 0; i < b.N; i++ {
		bb := bytes.NewBuffer(recordBytes)
		_, err := DeserializePrimitiveTestRecord(bb)
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
	someRecord, err := goavro.NewRecord(goavro.RecordSchema(string(schemaJson)))
	if err != nil {
		b.Fatal(err)
	}
	buf := new(bytes.Buffer)
	someRecord.Set("IntField", int32(1))
	someRecord.Set("LongField", int64(2))
	someRecord.Set("FloatField", float32(3.4))
	someRecord.Set("DoubleField", float64(5.6))
	someRecord.Set("StringField", "789")
	someRecord.Set("BoolField", true)
	someRecord.Set("BytesField", []byte{1, 2, 3, 4})

	err = codec.Encode(buf, someRecord)
	if err != nil {
		b.Fatal(err)
	}
	recordBytes := buf.Bytes()
	for i := 0; i < b.N; i++ {
		bb := bytes.NewBuffer(recordBytes)
		_, err := codec.Decode(bb)
		if err != nil {
			b.Fatal(err)
		}
	}

}
