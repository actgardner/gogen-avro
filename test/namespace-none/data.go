// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     namespace.avsc
 */
package avro

import (
	"github.com/actgardner/gogen-avro/compiler"
	"github.com/actgardner/gogen-avro/schema/canonical"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/vm/types"
	"io"
)

var DataUID []byte

func init() {
	t := NewData()
	DataUID = canonical.AvroCalcSchemaUID(t.Schema())
}

// Common information related to the event which must be included in any clean event
type Data struct {

	// Unique identifier for the event used for de-duplication and tracing.

	Uuid *UnionNullUUID

	// Fully qualified name of the host that generated the event that generated the data.

	Hostname *UnionNullString

	// Trace information not redundant with this object

	Trace *UnionNullTrace
}

func NewData() *Data {
	return &Data{}
}

func DeserializeData(r io.Reader) (*Data, error) {
	t := NewData()
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

func DeserializeDataFromSchema(r io.Reader, schema string) (*Data, error) {
	t := NewData()

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

func writeData(r *Data, w io.Writer) error {
	var err error

	err = writeUnionNullUUID(r.Uuid, w)
	if err != nil {
		return err
	}

	err = writeUnionNullString(r.Hostname, w)
	if err != nil {
		return err
	}

	err = writeUnionNullTrace(r.Trace, w)
	if err != nil {
		return err
	}

	return err
}

func (r *Data) Serialize(w io.Writer) error {
	return writeData(r, w)
}

func (r *Data) Schema() string {
	return "{\"doc\":\"Common information related to the event which must be included in any clean event\",\"fields\":[{\"default\":null,\"doc\":\"Unique identifier for the event used for de-duplication and tracing.\",\"name\":\"uuid\",\"type\":[\"null\",{\"doc\":\"A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014\",\"fields\":[{\"default\":\"\",\"name\":\"uuid\",\"type\":\"string\"}],\"name\":\"UUID\",\"namespace\":\"bodyworks.datatype\",\"type\":\"record\"}]},{\"default\":null,\"doc\":\"Fully qualified name of the host that generated the event that generated the data.\",\"name\":\"hostname\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"doc\":\"Trace information not redundant with this object\",\"name\":\"trace\",\"type\":[\"null\",{\"doc\":\"Trace\",\"fields\":[{\"default\":null,\"doc\":\"Trace Identifier\",\"name\":\"traceId\",\"type\":[\"null\",{\"doc\":\"A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014\",\"fields\":[{\"default\":\"\",\"name\":\"uuid\",\"type\":\"string\"}],\"name\":\"UUID\",\"namespace\":\"headerworks.datatype\",\"type\":\"record\"}]}],\"name\":\"Trace\",\"type\":\"record\"}]}],\"name\":\"bodyworks.Data\",\"type\":\"record\"}"
}

func (r *Data) SchemaName() string {
	return "bodyworks.Data"
}

func (_ *Data) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Data) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Data) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Data) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Data) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Data) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Data) SetString(v string)   { panic("Unsupported operation") }
func (_ *Data) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Data) Get(i int) types.Field {
	switch i {

	case 0:

		r.Uuid = NewUnionNullUUID()

		return r.Uuid

	case 1:

		r.Hostname = NewUnionNullString()

		return r.Hostname

	case 2:

		r.Trace = NewUnionNullTrace()

		return r.Trace

	}
	panic("Unknown field index")
}

func (r *Data) SetDefault(i int) {
	switch i {

	case 0:
		r.Uuid = NewUnionNullUUID()

		return

	case 1:
		r.Hostname = NewUnionNullString()

		return

	case 2:
		r.Trace = NewUnionNullTrace()

		return

	}
	panic("Unknown field index")
}

func (_ *Data) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Data) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Data) Finalize()                        {}
