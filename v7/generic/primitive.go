package generic

import (
	"github.com/actgardner/gogen-avro/v7/vm/types"
)

type primitiveDatum struct {
	value interface{}
}

func (r *primitiveDatum) Datum() interface{} {
	return r.value
}

func (r *primitiveDatum) SetBoolean(v bool) {
	r.value = v
}

func (r *primitiveDatum) SetInt(v int32) {
	r.value = v
}
func (r *primitiveDatum) SetLong(v int64) {
	r.value = v
}
func (r *primitiveDatum) SetFloat(v float32) {
	r.value = v
}
func (r *primitiveDatum) SetDouble(v float64) {
	r.value = v
}
func (r *primitiveDatum) SetBytes(v []byte) {
	r.value = v
}
func (r *primitiveDatum) SetString(v string) {
	r.value = v
}
func (r *primitiveDatum) Get(i int) types.Field { panic("") }
func (r *primitiveDatum) SetDefault(i int)      {}

func (r *primitiveDatum) AppendMap(key string) types.Field { panic("") }

func (r *primitiveDatum) AppendArray() types.Field { panic("") }

func (r *primitiveDatum) NullField(t int) {}
func (r *primitiveDatum) Finalize()       {}
