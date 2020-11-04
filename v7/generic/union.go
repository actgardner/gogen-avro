package generic

import (
	"github.com/actgardner/gogen-avro/v7/schema"
	"github.com/actgardner/gogen-avro/v7/vm/types"
)

type unionDatum struct {
	itemTypes []schema.AvroType
	datum     Datum
}

func (r *unionDatum) Datum() interface{} {
	return r.datum.Datum()
}

func (r *unionDatum) SetBoolean(v bool)   { panic("") }
func (r *unionDatum) SetInt(v int32)      { panic("") }
func (r *unionDatum) SetLong(v int64)     {}
func (r *unionDatum) SetFloat(v float32)  { panic("") }
func (r *unionDatum) SetDouble(v float64) { panic("") }
func (r *unionDatum) SetBytes(v []byte)   { panic("") }
func (r *unionDatum) SetString(v string)  { panic("") }
func (r *unionDatum) Get(i int) types.Field {
	r.datum = DatumForType(r.itemTypes[i])
	return r.datum
}
func (r *unionDatum) SetDefault(i int)                 {}
func (r *unionDatum) AppendMap(key string) types.Field { panic("") }
func (r *unionDatum) AppendArray() types.Field         { panic("") }
func (r *unionDatum) NullField(t int)                  { panic("") }
func (r *unionDatum) Finalize()                        {}
