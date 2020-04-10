package avro

import (
	"reflect"
	"testing"
)

// TestSchemaTag tests if the generated struct has the defined tag
func TestSchemaTag(t *testing.T) {
	tag := &StructTag{}
	field, ok := reflect.TypeOf(tag).Elem().FieldByName("ProductName")
	if !ok {
		t.Error("Field not found")
	}

	_, ok = field.Tag.Lookup("validate")
	if !ok {
		t.Error("Struct field has not the defined validate tag")
	}
}
