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

type UnionRecursiveUnionTestRecordTypeEnum int

const (
	UnionRecursiveUnionTestRecordTypeEnumRecursiveUnionTestRecord UnionRecursiveUnionTestRecordTypeEnum = 1
)

type UnionRecursiveUnionTestRecord struct {
	Null                     *types.NullVal
	RecursiveUnionTestRecord RecursiveUnionTestRecord
	UnionType                UnionRecursiveUnionTestRecordTypeEnum
}

func writeUnionRecursiveUnionTestRecord(r *UnionRecursiveUnionTestRecord, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionRecursiveUnionTestRecordTypeEnumRecursiveUnionTestRecord:
		return writeRecursiveUnionTestRecord(r.RecursiveUnionTestRecord, w)
	}
	return fmt.Errorf("invalid value for *UnionRecursiveUnionTestRecord")
}

func NewUnionRecursiveUnionTestRecord() *UnionRecursiveUnionTestRecord {
	return &UnionRecursiveUnionTestRecord{}
}

func (r *UnionRecursiveUnionTestRecord) Serialize(w io.Writer) error {
	return writeUnionRecursiveUnionTestRecord(r, w)
}

func DeserializeUnionRecursiveUnionTestRecord(r io.Reader) (*UnionRecursiveUnionTestRecord, error) {
	t := NewUnionRecursiveUnionTestRecord()
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

func DeserializeUnionRecursiveUnionTestRecordFromSchema(r io.Reader, schema string) (*UnionRecursiveUnionTestRecord, error) {
	t := NewUnionRecursiveUnionTestRecord()
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

func (r *UnionRecursiveUnionTestRecord) Schema() string {
	return "[\"null\",{\"fields\":[{\"name\":\"RecursiveField\",\"type\":[\"null\",\"RecursiveUnionTestRecord\"]}],\"name\":\"RecursiveUnionTestRecord\",\"type\":\"record\"}]"
}

func (_ *UnionRecursiveUnionTestRecord) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionRecursiveUnionTestRecord) SetLong(v int64) {

	r.UnionType = (UnionRecursiveUnionTestRecordTypeEnum)(v)
}

func (r *UnionRecursiveUnionTestRecord) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		r.RecursiveUnionTestRecord = NewRecursiveUnionTestRecord()
		return &types.Record{Target: (&r.RecursiveUnionTestRecord)}
	}
	panic("Unknown field index")
}
func (_ *UnionRecursiveUnionTestRecord) NullField(i int)  { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) SetDefault(i int) { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionRecursiveUnionTestRecord) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *UnionRecursiveUnionTestRecord) Finalize()                {}

func (r *UnionRecursiveUnionTestRecord) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionRecursiveUnionTestRecordTypeEnumRecursiveUnionTestRecord:
		return json.Marshal(map[string]interface{}{"RecursiveUnionTestRecord": r.RecursiveUnionTestRecord})
	}
	return nil, fmt.Errorf("invalid value for *UnionRecursiveUnionTestRecord")
}

func (r *UnionRecursiveUnionTestRecord) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["RecursiveUnionTestRecord"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.RecursiveUnionTestRecord)
	}
	return fmt.Errorf("invalid value for *UnionRecursiveUnionTestRecord")
}
