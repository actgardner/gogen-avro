// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     arrays.avsc
 */
package avro

import (
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
)

type UnionNullIntTypeEnum int

const (
	UnionNullIntTypeEnumNull UnionNullIntTypeEnum = 0

	UnionNullIntTypeEnumInt UnionNullIntTypeEnum = 1
)

type UnionNullInt struct {
	Null *types.NullVal

	Int int32

	UnionType UnionNullIntTypeEnum
}

func writeUnionNullInt(r *UnionNullInt, w io.Writer) error {
	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {

	case UnionNullIntTypeEnumNull:
		return vm.WriteNull(r.Null, w)

	case UnionNullIntTypeEnumInt:
		return vm.WriteInt(r.Int, w)

	}
	return fmt.Errorf("invalid value for *UnionNullInt")
}

func NewUnionNullInt() *UnionNullInt {
	return &UnionNullInt{}
}

func (_ *UnionNullInt) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionNullInt) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionNullInt) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionNullInt) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionNullInt) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionNullInt) SetString(v string)  { panic("Unsupported operation") }
func (r *UnionNullInt) SetLong(v int64) {
	r.UnionType = (UnionNullIntTypeEnum)(v)
}
func (r *UnionNullInt) Get(i int) types.Field {
	switch i {

	case 0:

		return r.Null

	case 1:

		return (*types.Int)(&r.Int)

	}
	panic("Unknown field index")
}
func (_ *UnionNullInt) SetDefault(i int)                 { panic("Unsupported operation") }
func (_ *UnionNullInt) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UnionNullInt) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *UnionNullInt) Finalize()                        {}
