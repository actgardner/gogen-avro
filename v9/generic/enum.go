package generic

import (
	"github.com/actgardner/gogen-avro/v9/schema"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)

type enumDatum struct {
	symbols []string
	value   string
}

func newEnumDatum(def *schema.EnumDefinition) *enumDatum {
	return &enumDatum{
		symbols: def.Symbols(),
	}
}

func (r *enumDatum) Datum() interface{} {
	return r.value
}

func (r *enumDatum) SetBoolean(v bool) { panic("cannot SetBoolean on generic enum") }
func (r *enumDatum) SetInt(v int32) {
	r.value = r.symbols[v]
}
func (r *enumDatum) SetLong(v int64)                  { panic("cannot SetLong on generic enum") }
func (r *enumDatum) SetFloat(v float32)               { panic("cannot SetFloat on generic enum") }
func (r *enumDatum) SetDouble(v float64)              { panic("cannot SetDouble on generic enum") }
func (r *enumDatum) SetBytes(v []byte)                { panic("cannot SetBytes on generic enum") }
func (r *enumDatum) SetString(v string)               { panic("cannot SetString on generic enum") }
func (r *enumDatum) Get(i int) types.Field            { panic("cannot Get on generic enum") }
func (r *enumDatum) SetDefault(i int)                 {}
func (r *enumDatum) AppendMap(key string) types.Field { panic("cannot AppendMap on generic enum") }
func (r *enumDatum) AppendArray() types.Field         { panic("cannot AppendArray on generic enum") }
func (r *enumDatum) NullField(t int)                  { panic("cannot NullField on generic enum") }
func (r *enumDatum) HintSize(t int)                   { panic("cannot HintSize on generic enum") }
func (r *enumDatum) Finalize()                        {}
