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
  {
    "NumberField": {
      "IntField": 1, 
      "LongField": 2, 
      "FloatField": 3.4, 
      "DoubleField": 5.6
    }, 
    "OtherField": {
      "StringField": "789", 
      "BoolField": true, 
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  },
  {
    "NumberField": {
      "IntField": 2147483647, 
      "LongField": 9223372036854775807, 
      "FloatField": 3.402823e+38, 
      "DoubleField": 1.7976931348623157e+308
    }, 
    "OtherField": {
      "StringField": "abcdghejw", 
      "BoolField": true, 
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  },
  {
    "NumberField": {
      "IntField": -2147483647, 
      "LongField": -9223372036854775807, 
      "FloatField": 3.402823e-38, 
      "DoubleField": 2.2250738585072014e-308
    }, 
    "OtherField": {
      "StringField": "jdnwjkendwedddedee", 
      "BoolField": true, 
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  }
]
`

func TestNestedFixture(t *testing.T) {
	fixtures := make([]NestedTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("nested.avsc")
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
		datum, _, err := codec.NativeFromBinary(buf.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		record := datum.(map[string]interface{})
		value := reflect.ValueOf(f)
		for i := 0; i < value.NumField(); i++ {
			fieldName := value.Type().Field(i).Name
			structVal := reflect.Indirect(value.Field(i))
			for j := 0; j < structVal.NumField(); j++ {
				nestedFieldName := structVal.Type().Field(j).Name
				avroVal, ok := record[fieldName]
				if !ok {
					t.Fatalf("GOT: %#v; WANT: %#v", ok, true)
				}
				nestedAvroVal, ok := avroVal.(map[string]interface{})[nestedFieldName]
				if !ok {
					t.Fatalf("GOT: %#v; WANT: %#v", ok, true)
				}
				nestedStructVal := structVal.Field(j).Interface()
				if !reflect.DeepEqual(nestedStructVal, nestedAvroVal) {
					t.Fatalf("Field %v not equal: %v != %v", fieldName+"."+nestedFieldName, nestedStructVal, nestedAvroVal)
				}
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]NestedTestRecord, 0)
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
		datum, err := DeserializeNestedTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
