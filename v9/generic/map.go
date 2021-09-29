package generic

import (
	"github.com/actgardner/gogen-avro/v9/schema"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)

type mapDatum struct {
	itemType schema.AvroType
	items    map[string]Datum
}

func (r *mapDatum) Datum() interface{} {
	v := make(map[string]interface{})
	for k, item := range r.items {
		v[k] = item.Datum()
	}
	return v
}

func (r *mapDatum) SetBoolean(v bool) {}

func (r *mapDatum) SetInt(v int32)      {}
func (r *mapDatum) SetLong(v int64)     {}
func (r *mapDatum) SetFloat(v float32)  {}
func (r *mapDatum) SetDouble(v float64) {}
func (r *mapDatum) SetBytes(v []byte)   {}
func (r *mapDatum) SetString(v string)  {}

func (r *mapDatum) Get(i int) types.Field { panic("cannot Get on generic map") }
func (r *mapDatum) SetDefault(i int)      {}

func (r *mapDatum) AppendMap(key string) types.Field {
	if r.items == nil {
		r.items = make(map[string]Datum)
	}
	d := DatumForType(r.itemType)
	r.items[key] = d
	return d
}

func (r *mapDatum) AppendArray() types.Field { panic("cannot AppendArray on generic map") }

func (r *mapDatum) NullField(t int) {}
func (r *mapDatum) Finalize()       {}
