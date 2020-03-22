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

var BodyworksDataUID []byte

func init() {
	t := NewBodyworksData()
	BodyworksDataUID = canonical.AvroCalcSchemaUID(t.Schema())
}

// Common information related to the event which must be included in any clean event
type BodyworksData struct {

	// Unique identifier for the event used for de-duplication and tracing.

	Uuid *UnionNullDatatypeUUID

	// Fully qualified name of the host that generated the event that generated the data.

	Hostname *UnionNullString

	// Trace information not redundant with this object

	Trace *UnionNullBodyworksTrace
}

func NewBodyworksData() *BodyworksData {
	return &BodyworksData{}
}

func DeserializeBodyworksData(r io.Reader) (*BodyworksData, error) {
	t := NewBodyworksData()
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

func DeserializeBodyworksDataFromSchema(r io.Reader, schema string) (*BodyworksData, error) {
	t := NewBodyworksData()
	err := canonical.AvroConsumeHeader(r)
	if err != nil {
		return nil, err
	}

	var deser *vm.Program
	deser, err = compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func writeBodyworksData(r *BodyworksData, w io.Writer) error {
	var err error

	err = writeUnionNullDatatypeUUID(r.Uuid, w)
	if err != nil {
		return err
	}

	err = writeUnionNullString(r.Hostname, w)
	if err != nil {
		return err
	}

	err = writeUnionNullBodyworksTrace(r.Trace, w)
	if err != nil {
		return err
	}

	return err
}

func (r *BodyworksData) Serialize(w io.Writer) error {
	return writeBodyworksData(r, w)
}

func (r *BodyworksData) Schema() string {
	return "{\"doc\":\"Common information related to the event which must be included in any clean event\",\"fields\":[{\"default\":null,\"doc\":\"Unique identifier for the event used for de-duplication and tracing.\",\"name\":\"uuid\",\"type\":[\"null\",{\"doc\":\"A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014\",\"fields\":[{\"default\":\"\",\"name\":\"uuid\",\"type\":\"string\"}],\"name\":\"UUID\",\"namespace\":\"bodyworks.datatype\",\"type\":\"record\"}]},{\"default\":null,\"doc\":\"Fully qualified name of the host that generated the event that generated the data.\",\"name\":\"hostname\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"doc\":\"Trace information not redundant with this object\",\"name\":\"trace\",\"type\":[\"null\",{\"doc\":\"Trace\",\"fields\":[{\"default\":null,\"doc\":\"Trace Identifier\",\"name\":\"traceId\",\"type\":[\"null\",{\"doc\":\"A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014\",\"fields\":[{\"default\":\"\",\"name\":\"uuid\",\"type\":\"string\"}],\"name\":\"UUID\",\"namespace\":\"headerworks.datatype\",\"type\":\"record\"}]}],\"name\":\"Trace\",\"type\":\"record\"}]}],\"name\":\"bodyworks.Data\",\"type\":\"record\"}"
}

func (r *BodyworksData) SchemaName() string {
	return "bodyworks.Data"
}

func (_ *BodyworksData) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *BodyworksData) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *BodyworksData) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *BodyworksData) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *BodyworksData) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *BodyworksData) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *BodyworksData) SetString(v string)   { panic("Unsupported operation") }
func (_ *BodyworksData) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *BodyworksData) Get(i int) types.Field {
	switch i {

	case 0:

		r.Uuid = NewUnionNullDatatypeUUID()

		return r.Uuid

	case 1:

		r.Hostname = NewUnionNullString()

		return r.Hostname

	case 2:

		r.Trace = NewUnionNullBodyworksTrace()

		return r.Trace

	}
	panic("Unknown field index")
}

func (r *BodyworksData) SetDefault(i int) {
	switch i {

	case 0:
		r.Uuid = NewUnionNullDatatypeUUID()

		return

	case 1:
		r.Hostname = NewUnionNullString()

		return

	case 2:
		r.Trace = NewUnionNullBodyworksTrace()

		return

	}
	panic("Unknown field index")
}

func (_ *BodyworksData) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *BodyworksData) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *BodyworksData) Finalize()                        {}
