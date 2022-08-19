// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     block.avsc
 *     header.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
	"io"
)

func writeMapBytes(r map[string]Bytes, w io.Writer) error {
	err := vm.WriteLong(int64(len(r)), w)
	if err != nil || len(r) == 0 {
		return err
	}
	for k, e := range r {
		err = vm.WriteString(k, w)
		if err != nil {
			return err
		}
		err = vm.WriteBytes(e, w)
		if err != nil {
			return err
		}
	}
	return vm.WriteLong(0, w)
}

type MapBytesWrapper struct {
	Target *map[string]Bytes
	keys   []string
	values []Bytes
}

func (_ *MapBytesWrapper) SetBoolean(v bool)     { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetInt(v int32)        { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetLong(v int64)       { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetFloat(v float32)    { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetDouble(v float64)   { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetBytes(v []byte)     { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetString(v string)    { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetUnionElem(v int64)  { panic("Unsupported operation") }
func (_ *MapBytesWrapper) Get(i int) types.Field { panic("Unsupported operation") }
func (_ *MapBytesWrapper) SetDefault(i int)      { panic("Unsupported operation") }

func (r *MapBytesWrapper) HintSize(s int) {
	if r.keys == nil {
		r.keys = make([]string, 0, s)
		r.values = make([]Bytes, 0, s)
	}
}

func (r *MapBytesWrapper) NullField(_ int) {
	panic("Unsupported operation")
}

func (r *MapBytesWrapper) Finalize() {
	for i := range r.keys {
		(*r.Target)[r.keys[i]] = r.values[i]
	}
}

func (r *MapBytesWrapper) AppendMap(key string) types.Field {
	r.keys = append(r.keys, key)
	var v Bytes
	r.values = append(r.values, v)
	return &BytesWrapper{Target: &r.values[len(r.values)-1]}
}

func (_ *MapBytesWrapper) AppendArray() types.Field { panic("Unsupported operation") }