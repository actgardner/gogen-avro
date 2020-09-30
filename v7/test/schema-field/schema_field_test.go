package avro

import (
	"testing"
)

// just to check compilation if the schema contains a field with name 'schema'
func TestCompilation(t *testing.T) {

	record := SchemaField{}
	val := record.AvroRecordSchema()

	if val == "" {
		t.Error()
	}

}
