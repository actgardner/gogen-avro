package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
	{"UnionField": null},
	{"UnionField":{"array":[1,2,3]}},
	{"UnionField":{"map":{"a":1, "b":3, "c": 5}}},
	{"UnionField":{"NestedUnionRecord":{"IntField":789}}}
]
`

func TestPrimitiveUnionFixture(t *testing.T) {
	fixtures := make([]ComplexUnionTestRecord, 0)
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

		datum, _, err := codec.NativeFromBinary(buf.Bytes())
		assert.Nil(t, err)

		record := datum.(map[string]interface{})
		recordField, ok := record["UnionField"]
		assert.Equal(t, true, ok)

		if f.UnionField == nil {
			if recordField != nil {
				t.Fatalf("Expected nil value")
			}
			continue
		}

		switch f.UnionField.UnionType {
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumArrayInt:
			arr := recordField.(map[string]interface{})["array"].([]interface{})
			for i, v := range arr {
				if v.(int32) != f.UnionField.ArrayInt[i] {
					t.Fatalf("Expected int value %v for union field, got %v", f.UnionField.ArrayInt[i], v)
				}
			}
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumMapInt:
			m := recordField.(map[string]interface{})["map"].(map[string]interface{})
			for k, v := range m {
				if v.(int32) != f.UnionField.MapInt[k] {
					t.Fatalf("Expected int value %v for union map key %v field, got %v", f.UnionField.MapInt[k], k, v)
				}
			}
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNestedUnionRecord:
			v, ok := recordField.(map[string]interface{})["NestedUnionRecord"].(map[string]interface{})["IntField"]
			if !ok {
				t.Fatalf("GOT: %#v; WANT: %#v", ok, true)
			}
			if v.(int32) != f.UnionField.NestedUnionRecord.IntField {
				t.Fatalf("Expected int value %v for nested record in union, got %v", f.UnionField.NestedUnionRecord.IntField, v)
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]ComplexUnionTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeComplexUnionTestRecord(&buf)
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}
