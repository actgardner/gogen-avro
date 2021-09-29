package generic

import (
	"github.com/actgardner/gogen-avro/v9/schema"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)

type unionDatum struct {
	itemTypes []schema.AvroType
	datum     Datum
}

func (r *unionDatum) Datum() interface{} {
	return r.datum.Datum()
}

func (r *unionDatum) SetBoolean(v bool)   { panic("cannot SetBoolean on generic union") }
func (r *unionDatum) SetInt(v int32)      { panic("cannot SetInt on generic union") }
func (r *unionDatum) SetLong(v int64)     {}
func (r *unionDatum) SetFloat(v float32)  { panic("cannot SetFloat on generic union") }
func (r *unionDatum) SetDouble(v float64) { panic("cannot SetDouble on generic union") }
func (r *unionDatum) SetBytes(v []byte)   { panic("cannot SetBytes on generic union") }
func (r *unionDatum) SetString(v string)  { panic("cannot SetString on generic union") }
func (r *unionDatum) Get(i int) types.Field {
	r.datum = DatumForType(r.itemTypes[i])
	return r.datum
}
func (r *unionDatum) SetDefault(i int)                 {}
func (r *unionDatum) AppendMap(key string) types.Field { panic("cannot AppendMap on generic union") }
func (r *unionDatum) AppendArray() types.Field         { panic("cannot AppendArray on generic union") }
func (r *unionDatum) NullField(t int)                  { panic("cannot NullField on generic union") }
func (r *unionDatum) Finalize()                        {}
