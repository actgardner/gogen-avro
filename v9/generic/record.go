package generic

import (
	"github.com/actgardner/gogen-avro/v9/schema"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)

type recordDatum struct {
	def    *schema.RecordDefinition
	fields []Datum
}

func newRecordDatum(def *schema.RecordDefinition) *recordDatum {
	return &recordDatum{
		def:    def,
		fields: make([]Datum, len(def.Fields())),
	}
}

func (r *recordDatum) Datum() interface{} {
	m := make(map[string]interface{})
	for i, f := range r.def.Fields() {
		m[f.Name()] = r.fields[i].Datum()
	}
	return m
}

func (r *recordDatum) SetBoolean(v bool)   { panic("cannot SetBoolean on generic record") }
func (r *recordDatum) SetInt(v int32)      { panic("cannot SetInt on generic record") }
func (r *recordDatum) SetLong(v int64)     { panic("cannot SetLong on generic record") }
func (r *recordDatum) SetFloat(v float32)  { panic("cannot SetFloat on generic record") }
func (r *recordDatum) SetDouble(v float64) { panic("cannot SetDouble on generic record") }
func (r *recordDatum) SetBytes(v []byte)   { panic("cannot SetBytes on generic record") }
func (r *recordDatum) SetString(v string)  { panic("cannot SetString on generic record") }
func (r *recordDatum) Get(i int) types.Field {
	field := r.def.Fields()[i]
	r.fields[i] = DatumForType(field.Type())
	return r.fields[i]
}
func (r *recordDatum) SetDefault(i int) {
	field := r.def.Fields()[i]
	r.fields[i] = &primitiveDatum{field.Default()}
}
func (r *recordDatum) AppendMap(key string) types.Field { panic("cannot AppendMap on generic record") }
func (r *recordDatum) AppendArray() types.Field         { panic("cannot AppendArray on generic record") }
func (r *recordDatum) NullField(t int) {
	r.fields[t] = &primitiveDatum{nil}
}
func (r *recordDatum) Finalize() {}
