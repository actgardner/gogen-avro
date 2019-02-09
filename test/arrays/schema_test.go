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
{"IntField": [1, -2147483647, 2147483647], "LongField": [2, 9223372036854775807, -9223372036854775807], "FloatField": [3.4, 3.402823e-38, 3.402823e+38], "DoubleField": [ 5.6, 2.2250738585072014e-308], "StringField": ["short", "789", "longer", "a slightly longer string"], "BoolField": [true, false], "BytesField": ["VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"]},
{"IntField":[], "LongField": [2], "FloatField": [], "DoubleField": [5.6], "StringField": [], "BoolField": [true], "BytesField": []},
{"IntField":[], "LongField": [], "FloatField": [], "DoubleField": [], "StringField": [], "BoolField": [], "BytesField": []}
]
`

func TestArrayFixture(t *testing.T) {
	fixtures := make([]ArrayTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	schemaJson, err := ioutil.ReadFile("arrays.avsc")
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
		record := datum.(map[string]interface{})

		value := reflect.ValueOf(f)
		for i := 0; i < value.NumField(); i++ {
			fieldName := value.Type().Field(i).Name
			avroVal, ok := record[fieldName]

			assert.Equal(t, true, ok)

			avroArray := avroVal.([]interface{})
			assert.Equal(t, len(avroArray), value.Field(i).Len())

			for j := 0; j < value.Field(i).Len(); j++ {
				avroArrayVal := avroArray[j]
				structArrayVal := value.Field(i).Index(j).Interface()
				assert.Equal(t, avroArrayVal, structArrayVal)
			}
		}
	}
}

func BenchmarkArrayRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := ArrayTestRecord{[]int32{1, 2, 3}, []int64{4, 5, 6}, []float64{3.4, 5.6, 7.8}, []string{"abc", "def", "ghi"}, []float32{10.1, 10.2, 10.3}, []bool{true, false}, [][]byte{{1, 2, 3, 4}}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := record.Serialize(buf)
		assert.Nil(b, err)
	}
}

func BenchmarkArrayGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("arrays.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"IntField":    []interface{}{int32(1), int32(2), int32(3)},
		"LongField":   []interface{}{int64(4), int64(5), int64(6)},
		"FloatField":  []interface{}{float32(10.1), float32(10.2), float32(10.3)},
		"DoubleField": []interface{}{float64(3.4), float64(5.6), float64(7.8)},
		"StringField": []interface{}{"abc", "def", "ghi"},
		"BoolField":   []interface{}{true, false},
		"BytesField":  []interface{}{[]byte{1, 2, 3, 4}},
	}
	buf := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := codec.BinaryFromNative(buf, someRecord)
		assert.Nil(b, err)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]ArrayTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeArrayTestRecord(&buf)
		assert.Nil(t, err)

		assert.Equal(t, *datum, f)
	}
}
