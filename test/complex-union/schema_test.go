package avro

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
{"UnionField":{"UnionType":0}},
{"UnionField":{"ArrayInt":[1,2,3], "UnionType":1}},
{"UnionField":{"MapInt":{"a":1, "b":3, "c": 5}, "UnionType":2}},
{"UnionField":{"NestedUnionRecord":{"IntField":789}, "UnionType":3}}
]
`

func TestPrimitiveUnionFixture(t *testing.T) {
	fixtures := make([]ComplexUnionTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("union.avsc")
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
		recordField, err := record.Get("UnionField")
		if err != nil {
			t.Fatal(err)
		}
		switch f.UnionField.UnionType {
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNull:
			if recordField != nil {
				t.Fatalf("Expected nil value for union field, got %v", recordField)
			}
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumArrayInt:
			arr := recordField.([]interface{})
			for i, v := range arr {
				if v.(int32) != f.UnionField.ArrayInt[i] {
					t.Fatalf("Expected int value %v for union field, got %v", f.UnionField.ArrayInt[i], v)
				}
			}
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumMapInt:
			m := recordField.(map[string]interface{})
			for k, v := range m {
				if v.(int32) != f.UnionField.MapInt[k] {
					t.Fatalf("Expected int value %v for union map key %v field, got %v", f.UnionField.MapInt[k], k, v)
				}
			}
		case UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNestedUnionRecord:
			v, err := recordField.(*goavro.Record).Get("IntField")
			if err != nil {
				t.Fatal(err)
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
		datum, err := DeserializeComplexUnionTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
