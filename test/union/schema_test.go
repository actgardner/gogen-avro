package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
{"UnionField":{"Int":1, "UnionType":0}},
{"UnionField":{"Long":2, "UnionType":1}},
{"UnionField":{"Float":3.4, "UnionType":2}},
{"UnionField":{"Double":5.6, "UnionType":3}},
{"UnionField":{"String":"testString", "UnionType":4}},
{"UnionField":{"Bool":true, "UnionType":5}},
{"UnionField":{"UnionType":6}}
]
`

func TestPrimitiveUnionFixture(t *testing.T) {
	fixtures := make([]PrimitiveUnionTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	schemaJson, err := ioutil.ReadFile("union.avsc")
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

		if got, want := len(remaining), 0; got != want {
			t.Fatalf("GOT: %#v; WANT: %#v", got, want)
		}
		record := datum.(map[string]interface{})
		value := reflect.ValueOf(f)

		for i := 0; i < value.NumField(); i++ {
			fieldName := value.Type().Field(i).Name
			fieldUnionIndex := int(value.Field(i).Elem().FieldByName("UnionType").Int())
			structVal := value.Field(i).Elem().Field(fieldUnionIndex).Interface()
			var avroVal interface{}
			top, ok := record[fieldName].(map[string]interface{})
			if ok {
				for _, v := range top {
					avroVal = v
					break
				}
			}
			switch structVal.(type) {
			case *types.NullVal:
				if avroVal != nil {
					t.Fatalf("Field %v not equal: %t != %t", fieldName, structVal, avroVal)
				}
			default:
				if !reflect.DeepEqual(structVal, avroVal) {
					t.Fatalf("Field %v not equal: %t != %t", fieldName, structVal, avroVal)
				}
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]PrimitiveUnionTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializePrimitiveUnionTestRecord(&buf)
		assert.Nil(t, err)

		assert.Equal(t, *datum, f)
	}
}
