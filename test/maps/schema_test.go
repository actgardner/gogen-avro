package avro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/linkedin/goavro"
	"io/ioutil"
	"reflect"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
{"IntField": {"small": 1, "min":-2147483647, "max":2147483647}, "LongField": {"small": 2, "min": 9223372036854775807, "max": -9223372036854775807}, "FloatField": {"small": 3.4, "verysmall": 3.402823e-38, "large": 3.402823e+38}, "DoubleField": {"small": 5.6, "verysmall": 2.2250738585072014e-308}, "StringField": {"short": "789", "longer": "a slightly longer string"}, "BoolField": {"true": true, "false":false}, "BytesField": {"small": "", "longer": ""}}
]
`

func TestMapFixture(t *testing.T) {
	fixtures := make([]MapTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("maps.avsc")
	if err != nil {
		t.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for i, f := range fixtures {
		fmt.Printf("Serializing fixture %v\n", i)
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%v\n", buf.Bytes())
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
			avroMap := avroVal.(map[string]interface{})
			if len(avroMap) != value.Field(i).Len() {
				t.Fatalf("Got %v keys from goavro but expected %v", len(avroMap), value.Field(i).Len())
			}
			for _, k := range value.Field(i).MapKeys() {
				keyString := k.Interface().(string)
				avroMapVal := avroMap[keyString]
				structMapVal := value.Field(i).MapIndex(k).Interface()
				if !reflect.DeepEqual(avroMapVal, structMapVal) {
					t.Fatalf("Field %v key %v not equal: %v != %v", fieldName, k, avroMapVal, structMapVal)
				}
			}
		}
	}
}

/*
func BenchmarkPrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}
		record.Serialize(buf)
	}
}

func BenchmarkPrimitiveGoavro(b *testing.B) {
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
		someRecord.Set("IntField", 1)
		someRecord.Set("LongField", 2)
		someRecord.Set("FloatField", 3.4)
		someRecord.Set("DoubleField", 5.6)
		someRecord.Set("StringField", "789")
		someRecord.Set("BoolField", true)
		someRecord.Set("BytesField", []byte{1, 2, 3, 4})

		codec.Encode(buf, someRecord)
	}

}
*/
