// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCES:
 *     block.avsc
 *     header.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/schema/canonical"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

var AvroContainerHeaderUID []byte

func init() {
	t := NewAvroContainerHeader()
	AvroContainerHeaderUID = canonical.AvroCalcSchemaUID(t.Schema())
}

type AvroContainerHeader struct {
	Magic Magic

	Meta *MapBytes

	Sync Sync
}

func NewAvroContainerHeader() *AvroContainerHeader {
	return &AvroContainerHeader{}
}

func DeserializeAvroContainerHeader(r io.Reader) (*AvroContainerHeader, error) {
	t := NewAvroContainerHeader()
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

func DeserializeAvroContainerHeaderFromSchema(r io.Reader, schema string) (*AvroContainerHeader, error) {
	t := NewAvroContainerHeader()

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

func writeAvroContainerHeader(r *AvroContainerHeader, w io.Writer) error {
	var err error

	err = writeMagic(r.Magic, w)
	if err != nil {
		return err
	}

	err = writeMapBytes(r.Meta, w)
	if err != nil {
		return err
	}

	err = writeSync(r.Sync, w)
	if err != nil {
		return err
	}

	return err
}

func (r *AvroContainerHeader) Serialize(w io.Writer) error {
	return writeAvroContainerHeader(r, w)
}

func (r *AvroContainerHeader) Schema() string {
	return "{\"fields\":[{\"name\":\"magic\",\"type\":{\"name\":\"Magic\",\"size\":4,\"type\":\"fixed\"}},{\"name\":\"meta\",\"type\":{\"type\":\"map\",\"values\":\"bytes\"}},{\"name\":\"sync\",\"type\":{\"name\":\"Sync\",\"size\":16,\"type\":\"fixed\"}}],\"name\":\"AvroContainerHeader\",\"type\":\"record\"}"
}

func (r *AvroContainerHeader) SchemaName() string {
	return "AvroContainerHeader"
}

func (_ *AvroContainerHeader) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetString(v string)   { panic("Unsupported operation") }
func (_ *AvroContainerHeader) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *AvroContainerHeader) Get(i int) types.Field {
	switch i {

	case 0:

		return (*MagicWrapper)(&r.Magic)

	case 1:

		r.Meta = NewMapBytes()

		return r.Meta

	case 2:

		return (*SyncWrapper)(&r.Sync)

	}
	panic("Unknown field index")
}

func (r *AvroContainerHeader) SetDefault(i int) {
	switch i {

	}
	panic("Unknown field index")
}

func (_ *AvroContainerHeader) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *AvroContainerHeader) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *AvroContainerHeader) Finalize()                        {}
