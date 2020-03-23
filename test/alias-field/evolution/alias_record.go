// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     evolution.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

type AliasRecord struct {
	B string

	D string
}

var AliasRecordAvroCRC64Fingerprint = []byte{0x77, 0x54, 0xf0, 0x8a, 0x2b, 0xc9, 0xa8, 0xce}

func NewAliasRecord() *AliasRecord {
	return &AliasRecord{}
}

func DeserializeAliasRecord(r io.Reader) (*AliasRecord, error) {
	t := NewAliasRecord()
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

func DeserializeAliasRecordFromSchema(r io.Reader, schema string) (*AliasRecord, error) {
	t := NewAliasRecord()

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

func writeAliasRecord(r *AliasRecord, w io.Writer) error {
	var err error

	err = vm.WriteString(r.B, w)
	if err != nil {
		return err
	}

	err = vm.WriteString(r.D, w)
	if err != nil {
		return err
	}

	return err
}

func (r *AliasRecord) Serialize(w io.Writer) error {
	return writeAliasRecord(r, w)
}

func (r *AliasRecord) Schema() string {
	return "{\"fields\":[{\"aliases\":[\"a\"],\"name\":\"b\",\"type\":\"string\"},{\"name\":\"d\",\"type\":\"string\"}],\"name\":\"AliasRecord\",\"type\":\"record\"}"
}

func (r *AliasRecord) SchemaName() string {
	return "AliasRecord"
}

func (_ *AliasRecord) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *AliasRecord) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *AliasRecord) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *AliasRecord) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *AliasRecord) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *AliasRecord) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *AliasRecord) SetString(v string)   { panic("Unsupported operation") }
func (_ *AliasRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *AliasRecord) Get(i int) types.Field {
	switch i {

	case 0:

		return (*types.String)(&r.B)

	case 1:

		return (*types.String)(&r.D)

	}
	panic("Unknown field index")
}

func (r *AliasRecord) SetDefault(i int) {
	switch i {

	}
	panic("Unknown field index")
}

func (_ *AliasRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *AliasRecord) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *AliasRecord) Finalize()                        {}

func (_ *AliasRecord) AvroCRC64Fingerprint() []byte {
	return AliasRecordAvroCRC64Fingerprint
}
