package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
    {
        "IntField": {
            "small": 1, "min":-2147483647, "max":2147483647
        },
        "LongField": {
            "small": 2, "min": 9223372036854775807, "max": -9223372036854775807
        },
        "FloatField": {
            "small": 3.4, "verysmall": 3.402823e-38, "large": 3.402823e+38
        }, 
        "DoubleField": {
            "small": 5.6, "verysmall": 2.2250738585072014e-308
        }, 
        "StringField": {
            "short": "789", "longer": "a slightly longer string"
        },
        "BoolField": {
            "true": true, "false": false
        }, 
        "BytesField": {
            "small": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n", "longer": "VGhpcyBpcyBhIG11Y2ggbG9uZ2VyIHRlc3Qgc3RyaW5nIGxvbmcgbG9uZw=="
        }
    },
    {
        "IntField": {}, 
        "LongField": {},
        "FloatField": {},
        "DoubleField": {},
        "StringField": {},
        "BoolField": {"true": true},
        "BytesField": {}
    }
]
`

func BenchmarkMapRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := MapTestRecord{
		map[string]int32{"value1": 1, "value2": 2, "value3": 3},
		map[string]int64{"value1": 1, "value2": 2, "value3": 3},
		map[string]float64{"value1": 1, "value2": 2, "value3": 3},
		map[string]string{"value1": "12345", "value2": "67890", "value3": "abcdefg"},
		map[string]float32{"value1": 1, "value2": 2, "value3": 3},
		map[string]bool{"true": true, "false": false},
		map[string][]byte{"value1": {1, 2, 3, 4}, "value2": {100, 200, 255}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := record.Serialize(buf)
		assert.Nil(b, err)
	}
}

func BenchmarkMapGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("maps.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"IntField":    map[string]interface{}{"value1": int32(1), "value2": int32(2), "value3": int32(3)},
		"LongField":   map[string]interface{}{"value1": int64(1), "value2": int32(2), "value3": int32(3)},
		"FloatField":  map[string]interface{}{"value1": float32(1), "value2": float32(2), "value3": float32(3)},
		"DoubleField": map[string]interface{}{"value1": float32(1), "value2": float32(2), "value3": float32(3)},
		"StringField": map[string]interface{}{"value1": "12345", "value2": "67890", "value3": "abcdefg"},
		"BoolField":   map[string]interface{}{"true": true, "false": false},
		"BytesField":  map[string]interface{}{"value1": []byte{1, 2, 3, 4}, "value2": []byte{100, 200, 255}},
	}
	buf := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := codec.BinaryFromNative(buf, someRecord)
		assert.Nil(b, err)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]MapTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeMapTestRecord(&buf)
		assert.Nil(t, err)
		assert.Equal(t, datum.IntField, f.IntField)
		assert.Equal(t, datum.LongField, f.LongField)
		assert.Equal(t, datum.FloatField, f.FloatField)
		assert.Equal(t, datum.DoubleField, f.DoubleField)
		assert.Equal(t, datum.StringField, f.StringField)
		assert.Equal(t, datum.BoolField, f.BoolField)
		assert.Equal(t, datum.BytesField, f.BytesField)
	}
}
