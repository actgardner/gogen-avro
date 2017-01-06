package test

import (
	"bytes"
	"fmt"
	"io"
	"github.com/linkedin/goavro"
	"reflect"
	"testing"
)

type Serializable interface {
	Serialize(io.Writer) error
}

func RoundTripGoAvroTest(fixture Serializable, codec goavro.Codec, t *testing.T) {
        var buf bytes.Buffer
	err := fixture.Serialize(&buf)
	if err != nil {
		t.Fatal(err)
	}
	datum, err := codec.Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	compareGoAvroRecord(fixture, datum, t)
}

func compareGoAvroRecord(structVal interface{}, avroVal interface{}, t *testing.T) {
	fmt.Printf("Comparing records %v %v\n", structVal, avroVal)
	record := avroVal.(*goavro.Record)
	value := reflect.Indirect(reflect.ValueOf(structVal))
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldStructVal := value.Field(i).Interface()
		fieldAvroVal, err := record.Get(fieldName)
		if err != nil {
			t.Fatal(err)
		}
		compareGoAvroFields(fieldStructVal, fieldAvroVal, t)
	}
}

func compareGoAvroFields(structVal interface{}, avroVal interface{}, t *testing.T) {
	fmt.Printf("Comparing fields %v %v\n", structVal, avroVal)
	switch avroVal.(type) {
		case []interface{}: compareGoAvroArray(structVal, avroVal, t)
		case map[string]interface{}: compareGoAvroMap(structVal, avroVal, t)
		case *goavro.Record: compareGoAvroRecord(structVal, avroVal, t)
		default: compareGoAvroPrimitive(structVal, avroVal, t)
	}
}

func compareGoAvroPrimitive(structVal interface{}, avroVal interface{}, t *testing.T) {
	fmt.Printf("Comparing primitives %v %v\n", structVal, avroVal)
	if !reflect.DeepEqual(structVal, avroVal) {
		t.Fatalf("Fields not equal: %v != %v", structVal, avroVal)
	}
}

func compareGoAvroArray(structVal interface{}, avroVal interface{}, t *testing.T) {
	fmt.Printf("Comparing arrays %v %v\n", structVal, avroVal)
	arrayField := reflect.ValueOf(structVal)
	avroArray := avroVal.([]interface{})

	if len(avroArray) != arrayField.Len() {
		t.Fatalf("Got %v elements from goavro but expected %v", len(avroArray), arrayField.Len())
	}

	for j := 0; j < arrayField.Len(); j++ {
		avroArrayVal := avroArray[j]
		structArrayVal := arrayField.Index(j).Interface()
		compareGoAvroFields(structArrayVal, avroArrayVal, t)
	}
}

func compareGoAvroMap(structVal interface{}, avroVal interface{}, t *testing.T) {
	fmt.Printf("Comparing maps %v %v\n", structVal, avroVal)
	mapField := reflect.ValueOf(structVal)
	avroMap := avroVal.(map[string]interface{})

	if len(avroMap) != mapField.Len() {
		t.Fatalf("Got %v elements from goavro but expected %v", len(avroMap), mapField.Len())
	}

	for _, k := range mapField.MapKeys() {
		keyString := k.Interface().(string)
		avroMapVal := avroMap[keyString]
		structMapVal := mapField.MapIndex(k).Interface()
		compareGoAvroFields(structMapVal, avroMapVal, t)
	}
}
