package avro

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/linkedin/goavro.v1"
	"io/ioutil"
	"reflect"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
{"IntField": [1, -2147483647, 2147483647], "LongField": [2, 9223372036854775807, -9223372036854775807], "FloatField": [3.4, 3.402823e-38, 3.402823e+38], "DoubleField": [ 5.6, 2.2250738585072014e-308], "StringField": ["short", "789", "longer", "a slightly longer string"], "BoolField": [true, false], "BytesField": ["VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"]},
{"IntField":[], "LongField": [2], "FloatField": [], "DoubleField": [5.6], "StringField": [], "BoolField": [true], "BytesField": []},
{"IntField":[], "LongField": [], "FloatField": [], "DoubleField": [], "StringField": [], "BoolField": [], "BytesField": []}
]
`

func TestArrayFixture(t *testing.T) {
	fixtures := make([]ArrayTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("arrays.avsc")
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
		datum, err := codec.Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		record := datum.(*goavro.Record)
		value := reflect.ValueOf(f)
		for i := 0; i < value.NumField(); i++ {
			fieldName := value.Type().Field(i).Name
			avroVal, err := record.Get(fieldName)
			if err != nil {
				t.Fatal(err)
			}
			avroArray := avroVal.([]interface{})
			if len(avroArray) != value.Field(i).Len() {
				t.Fatalf("Got %v elements from goavro but expected %v", len(avroArray), value.Field(i).Len())
			}
			for j := 0; j < value.Field(i).Len(); j++ {
				avroArrayVal := avroArray[j]
				structArrayVal := value.Field(i).Index(j).Interface()
				if !reflect.DeepEqual(avroArrayVal, structArrayVal) {
					t.Fatalf("Field %v element %v not equal: %v != %v", fieldName, j, avroArrayVal, structArrayVal)
				}
			}
		}
	}
}

func TestArrayDefaults(t *testing.T) {
	record := NewArrayTestRecord()
	assert.Equal(t, record.IntField, []int32{1, 2, 3, 4})
	assert.Equal(t, record.LongField, []int64{5, 6, 7, 8})
	assert.Equal(t, record.FloatField, []float32{1.23, 3.45})
	assert.Equal(t, record.DoubleField, []float64{1.5, 2.4})
	assert.Equal(t, record.StringField, []string{"abc", "def"})
	assert.Equal(t, record.BoolField, []bool{true, false})
	assert.Equal(t, record.BytesField, [][]byte{[]byte("abc"), []byte("def")})
}

func BenchmarkArrayRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		record := ArrayTestRecord{[]int32{1, 2, 3}, []int64{4, 5, 6}, []float64{3.4, 5.6, 7.8}, []string{"abc", "def", "ghi"}, []float32{10.1, 10.2, 10.3}, []bool{true, false}, [][]byte{{1, 2, 3, 4}}}
		err := record.Serialize(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArrayGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("arrays.avsc")
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
		err := someRecord.Set("IntField", []interface{}{int32(1), int32(2), int32(3)})
		if err != nil {
			b.Fatal(err)
		}
		err = someRecord.Set("LongField", []interface{}{int64(4), int64(5), int64(6)})
		if err != nil {
			b.Fatal(err)
		}
		err = someRecord.Set("FloatField", []interface{}{float32(10.1), float32(10.2), float32(10.3)})
		if err != nil {
			b.Fatal(err)
		}
		err = someRecord.Set("DoubleField", []interface{}{float64(3.4), float64(5.6), float64(7.8)})
		if err != nil {
			b.Fatal(err)
		}
		err = someRecord.Set("StringField", []interface{}{"abc", "def", "ghi"})
		if err != nil {
			b.Fatal(err)
		}
		err = someRecord.Set("BoolField", []interface{}{true, false})
		if err != nil {
			b.Fatal(err)
		}
		err = someRecord.Set("BytesField", []interface{}{[]byte{1, 2, 3, 4}})
		if err != nil {
			b.Fatal(err)
		}

		err = codec.Encode(buf, someRecord)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]ArrayTestRecord, 0)
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
		datum, err := DeserializeArrayTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
