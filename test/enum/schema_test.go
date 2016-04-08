package avro

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	"io/ioutil"
	"reflect"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
{"EnumField":{"Int":1, "UnionType":0}},
{"EnumField":{"Long":2, "UnionType":1}},
{"EnumField":{"Float":3.4, "UnionType":2}},
{"EnumField":{"Double":5.6, "UnionType":3}},
{"EnumField":{"String":"testString", "UnionType":4}},
{"EnumField":{"Bool":true, "UnionType":5}},
{"EnumField":{"Bytes":"VGhpcyBpcyBhIHRlc3Qgc3RyaW5n", "UnionType":6}},
{"EnumField":{"UnionType":7}}
]
`

func TestPrimitiveEnumFixture(t *testing.T) {
	fixtures := make([]PrimitiveEnumTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("enum.avsc")
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
			fieldUnionIndex := int(value.Field(i).FieldByName("UnionType").Int())
			structVal := value.Field(i).Field(fieldUnionIndex).Interface()
			avroVal, err := record.Get(fieldName)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(structVal, avroVal) {
				t.Fatalf("Field %v not equal: %v != %v", fieldName, structVal, avroVal)
			}
		}
	}
}
