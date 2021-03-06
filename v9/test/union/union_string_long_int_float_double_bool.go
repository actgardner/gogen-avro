// Code generated by github.com/actgardner/gogen-avro/v8. DO NOT EDIT.
/*
 * SOURCE:
 *     schema.avsc
 */
package avro

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v9/compiler"
	"github.com/actgardner/gogen-avro/v9/vm"
	"github.com/actgardner/gogen-avro/v9/vm/types"
)

type UnionStringLongIntFloatDoubleBoolTypeEnum int

const (
	UnionStringLongIntFloatDoubleBoolTypeEnumString UnionStringLongIntFloatDoubleBoolTypeEnum = 0

	UnionStringLongIntFloatDoubleBoolTypeEnumLong UnionStringLongIntFloatDoubleBoolTypeEnum = 1

	UnionStringLongIntFloatDoubleBoolTypeEnumInt UnionStringLongIntFloatDoubleBoolTypeEnum = 2

	UnionStringLongIntFloatDoubleBoolTypeEnumFloat UnionStringLongIntFloatDoubleBoolTypeEnum = 3

	UnionStringLongIntFloatDoubleBoolTypeEnumDouble UnionStringLongIntFloatDoubleBoolTypeEnum = 4

	UnionStringLongIntFloatDoubleBoolTypeEnumBool UnionStringLongIntFloatDoubleBoolTypeEnum = 6
)

type UnionStringLongIntFloatDoubleBool struct {
	String    string
	Long      int64
	Int       int32
	Float     float32
	Double    float64
	Null      *types.NullVal
	Bool      bool
	UnionType UnionStringLongIntFloatDoubleBoolTypeEnum
}

func writeUnionStringLongIntFloatDoubleBool(r *UnionStringLongIntFloatDoubleBool, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(5, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionStringLongIntFloatDoubleBoolTypeEnumString:
		return vm.WriteString(r.String, w)
	case UnionStringLongIntFloatDoubleBoolTypeEnumLong:
		return vm.WriteLong(r.Long, w)
	case UnionStringLongIntFloatDoubleBoolTypeEnumInt:
		return vm.WriteInt(r.Int, w)
	case UnionStringLongIntFloatDoubleBoolTypeEnumFloat:
		return vm.WriteFloat(r.Float, w)
	case UnionStringLongIntFloatDoubleBoolTypeEnumDouble:
		return vm.WriteDouble(r.Double, w)
	case UnionStringLongIntFloatDoubleBoolTypeEnumBool:
		return vm.WriteBool(r.Bool, w)
	}
	return fmt.Errorf("invalid value for *UnionStringLongIntFloatDoubleBool")
}

func NewUnionStringLongIntFloatDoubleBool() *UnionStringLongIntFloatDoubleBool {
	return &UnionStringLongIntFloatDoubleBool{}
}

func (r *UnionStringLongIntFloatDoubleBool) Serialize(w io.Writer) error {
	return writeUnionStringLongIntFloatDoubleBool(r, w)
}

func DeserializeUnionStringLongIntFloatDoubleBool(r io.Reader) (*UnionStringLongIntFloatDoubleBool, error) {
	t := NewUnionStringLongIntFloatDoubleBool()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func DeserializeUnionStringLongIntFloatDoubleBoolFromSchema(r io.Reader, schema string) (*UnionStringLongIntFloatDoubleBool, error) {
	t := NewUnionStringLongIntFloatDoubleBool()
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func (r *UnionStringLongIntFloatDoubleBool) Schema() string {
	return "[\"string\",\"long\",\"int\",\"float\",\"double\",\"null\",\"boolean\"]"
}

func (_ *UnionStringLongIntFloatDoubleBool) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionStringLongIntFloatDoubleBool) SetLong(v int64) {

	r.UnionType = (UnionStringLongIntFloatDoubleBoolTypeEnum)(v)
}

func (r *UnionStringLongIntFloatDoubleBool) Get(i int) types.Field {

	switch i {
	case 0:
		return &types.String{Target: (&r.String)}
	case 1:
		return &types.Long{Target: (&r.Long)}
	case 2:
		return &types.Int{Target: (&r.Int)}
	case 3:
		return &types.Float{Target: (&r.Float)}
	case 4:
		return &types.Double{Target: (&r.Double)}
	case 5:
		return r.Null
	case 6:
		return &types.Boolean{Target: (&r.Bool)}
	}
	panic("Unknown field index")
}
func (_ *UnionStringLongIntFloatDoubleBool) NullField(i int)  { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) SetDefault(i int) { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionStringLongIntFloatDoubleBool) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *UnionStringLongIntFloatDoubleBool) Finalize()                {}

func (r *UnionStringLongIntFloatDoubleBool) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionStringLongIntFloatDoubleBoolTypeEnumString:
		return json.Marshal(map[string]interface{}{"string": r.String})
	case UnionStringLongIntFloatDoubleBoolTypeEnumLong:
		return json.Marshal(map[string]interface{}{"long": r.Long})
	case UnionStringLongIntFloatDoubleBoolTypeEnumInt:
		return json.Marshal(map[string]interface{}{"int": r.Int})
	case UnionStringLongIntFloatDoubleBoolTypeEnumFloat:
		return json.Marshal(map[string]interface{}{"float": r.Float})
	case UnionStringLongIntFloatDoubleBoolTypeEnumDouble:
		return json.Marshal(map[string]interface{}{"double": r.Double})
	case UnionStringLongIntFloatDoubleBoolTypeEnumBool:
		return json.Marshal(map[string]interface{}{"boolean": r.Bool})
	}
	return nil, fmt.Errorf("invalid value for *UnionStringLongIntFloatDoubleBool")
}

func (r *UnionStringLongIntFloatDoubleBool) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["string"]; ok {
		r.UnionType = 0
		return json.Unmarshal([]byte(value), &r.String)
	}
	if value, ok := fields["long"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.Long)
	}
	if value, ok := fields["int"]; ok {
		r.UnionType = 2
		return json.Unmarshal([]byte(value), &r.Int)
	}
	if value, ok := fields["float"]; ok {
		r.UnionType = 3
		return json.Unmarshal([]byte(value), &r.Float)
	}
	if value, ok := fields["double"]; ok {
		r.UnionType = 4
		return json.Unmarshal([]byte(value), &r.Double)
	}
	if value, ok := fields["boolean"]; ok {
		r.UnionType = 6
		return json.Unmarshal([]byte(value), &r.Bool)
	}
	return fmt.Errorf("invalid value for *UnionStringLongIntFloatDoubleBool")
}
