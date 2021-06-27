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

func (r *enumDatum) SetBoolean(v bool) { panic("") }
func (r *enumDatum) SetInt(v int32)    { panic("") }
func (r *enumDatum) SetLong(v int64) {
	r.value = r.symbols[v]
}
func (r *enumDatum) SetFloat(v float32)               { panic("") }
func (r *enumDatum) SetDouble(v float64)              { panic("") }
func (r *enumDatum) SetBytes(v []byte)                { panic("") }
func (r *enumDatum) SetString(v string)               { panic("") }
func (r *enumDatum) Get(i int) types.Field            { panic("") }
func (r *enumDatum) SetDefault(i int)                 {}
func (r *enumDatum) AppendMap(key string) types.Field { panic("") }
func (r *enumDatum) AppendArray() types.Field         { panic("") }
func (r *enumDatum) NullField(t int)                  { panic("") }
func (r *enumDatum) Finalize()                        { panic("") }
