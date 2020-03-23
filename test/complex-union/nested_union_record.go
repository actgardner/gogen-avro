// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     union.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

type NestedUnionRecord struct {
	IntField int32
}

var NestedUnionRecordAvroCRC64Fingerprint = []byte{0x23, 0xb6, 0xed, 0xa0, 0x87, 0x46, 0x3d, 0xef}

func NewNestedUnionRecord() *NestedUnionRecord {
	return &NestedUnionRecord{}
}

func DeserializeNestedUnionRecord(r io.Reader) (*NestedUnionRecord, error) {
	t := NewNestedUnionRecord()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func DeserializeNestedUnionRecordFromSchema(r io.Reader, schema string) (*NestedUnionRecord, error) {
	t := NewNestedUnionRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func writeNestedUnionRecord(r *NestedUnionRecord, w io.Writer) error {
	var err error

	err = vm.WriteInt(r.IntField, w)
	if err != nil {
		return err
	}

	return err
}

func (r *NestedUnionRecord) Serialize(w io.Writer) error {
	return writeNestedUnionRecord(r, w)
}

func (r *NestedUnionRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"IntField\",\"type\":\"int\"}],\"name\":\"NestedUnionRecord\",\"type\":\"record\"}"
}

func (r *NestedUnionRecord) SchemaName() string {
	return "NestedUnionRecord"
}

func (_ *NestedUnionRecord) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetString(v string)   { panic("Unsupported operation") }
func (_ *NestedUnionRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *NestedUnionRecord) Get(i int) types.Field {
	switch i {

	case 0:

		return (*types.Int)(&r.IntField)

	}
	panic("Unknown field index")
}

func (r *NestedUnionRecord) SetDefault(i int) {
	switch i {

	}
	panic("Unknown field index")
}

func (_ *NestedUnionRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *NestedUnionRecord) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *NestedUnionRecord) Finalize()                        {}

func (_ *NestedUnionRecord) AvroCRC64Fingerprint() []byte {
	return NestedUnionRecordAvroCRC64Fingerprint
}
