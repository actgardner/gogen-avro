package zero_copy

import (
	avro "github.com/actgardner/gogen-avro/v8/schema"
)

type Record struct {
	FieldOffset []int
}

func (r *Record) SetOffset(field, offset int) {
	r.FieldOffset[field] = offset
}

func generateRecord(def *avro.RecordDefinition) string {
	return ""
}
