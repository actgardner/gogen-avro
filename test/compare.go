package avrotest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/linkedin/goavro"
	"io/ioutil"
	"reflect"
	"testing"
)

func compareRecord(t, r *goavro.Record, s interface{}) {
	value := reflect.ValueOf(s)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		structVal := value.Field(i).Interface()
		avroVal, err := record.Get(fieldName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(structVal, avroVal) {
			t.Fatalf("Field %v not equal: %v != %v", fieldName, structVal, avroVal)
		}
	}
}

/* Compare a map from goavro to a map in the struct */
func compareMap(t, rMap map[string]interface{}, m *reflect.Value) {
	if len(rMap) != m.Len() {
		t.Fatalf("Got %v keys from goavro but expected %v", len(rMap), m.Len())
	}
	for _, k := range m.MapKeys() {
		keyString := k.Interface().(string)
		rVal := rMap[keyString]
		mVal := m.MapIndex(k).Interface()
		if !reflect.DeepEqual(rVal, mVal) {
			t.Fatalf("Map key %v not equal: %v != %v", keyString, rVal, mVal)
		}
	}
}

func compareArray(t, r []interface{}, m *reflect.Value) {
	if len(r) != m.Len() {
		t.Fatalf("Got %v elements from goavro but expected %v", len(r), m.Len())
	}
	for i := 0; i < m.Len(); i++ {
		rVal := r[i]
		mVal := m.Index(i).Interface()
		if !reflect.DeepEqual(rVal, mVal) {
			t.Fatalf("Aray element %v not equal: %v != %v", i, rVal, mVal)
		}
	}
}
