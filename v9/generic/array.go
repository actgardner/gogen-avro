package generic

import (
	"github.com/actgardner/gogen-avro/v9/schema"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)

type arrayDatum struct {
	itemType schema.AvroType
	items    []Datum
}

func (r *arrayDatum) Datum() interface{} {
	v := make([]interface{}, len(r.items))
	for i, item := range r.items {
		v[i] = item.Datum()
	}
	return v
}

func (r *arrayDatum) SetBoolean(v bool) {}

func (r *arrayDatum) SetInt(v int32)      {}
func (r *arrayDatum) SetLong(v int64)     {}
func (r *arrayDatum) SetFloat(v float32)  {}
func (r *arrayDatum) SetDouble(v float64) {}
func (r *arrayDatum) SetBytes(v []byte)   {}
func (r *arrayDatum) SetString(v string)  {}

func (r *arrayDatum) Get(i int) types.Field { panic("cannot Get on generic array") }
func (r *arrayDatum) SetDefault(i int)      {}

func (r *arrayDatum) AppendMap(key string) types.Field { panic("cannot AppendMap on generic array") }

func (r *arrayDatum) AppendArray() types.Field {
	d := DatumForType(r.itemType)
	r.items = append(r.items, d)
	return d
}

func (r *arrayDatum) NullField(t int) {}
func (r *arrayDatum) Finalize()       {}
