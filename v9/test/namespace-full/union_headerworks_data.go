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

type UnionHeaderworksDataTypeEnum int

const (
	UnionHeaderworksDataTypeEnumHeaderworksData UnionHeaderworksDataTypeEnum = 1
)

type UnionHeaderworksData struct {
	Null            *types.NullVal
	HeaderworksData HeaderworksData
	UnionType       UnionHeaderworksDataTypeEnum
}

func writeUnionHeaderworksData(r *UnionHeaderworksData, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionHeaderworksDataTypeEnumHeaderworksData:
		return writeHeaderworksData(r.HeaderworksData, w)
	}
	return fmt.Errorf("invalid value for *UnionHeaderworksData")
}

func NewUnionHeaderworksData() *UnionHeaderworksData {
	return &UnionHeaderworksData{}
}

func (r *UnionHeaderworksData) Serialize(w io.Writer) error {
	return writeUnionHeaderworksData(r, w)
}

func DeserializeUnionHeaderworksData(r io.Reader) (*UnionHeaderworksData, error) {
	t := NewUnionHeaderworksData()
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

func DeserializeUnionHeaderworksDataFromSchema(r io.Reader, schema string) (*UnionHeaderworksData, error) {
	t := NewUnionHeaderworksData()
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

func (r *UnionHeaderworksData) Schema() string {
	return "[\"null\",{\"doc\":\"Common information related to the event which must be included in any clean event\",\"fields\":[{\"default\":null,\"doc\":\"Unique identifier for the event used for de-duplication and tracing.\",\"name\":\"uuid\",\"type\":[\"null\",{\"doc\":\"A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014\",\"fields\":[{\"default\":\"\",\"name\":\"uuid\",\"type\":\"string\"}],\"name\":\"UUID\",\"namespace\":\"headerworks.datatype\",\"type\":\"record\"}]},{\"default\":null,\"doc\":\"Fully qualified name of the host that generated the event that generated the data.\",\"name\":\"hostname\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"doc\":\"Trace information not redundant with this object\",\"name\":\"trace\",\"type\":[\"null\",{\"doc\":\"Trace\",\"fields\":[{\"default\":null,\"doc\":\"Trace Identifier\",\"name\":\"traceId\",\"type\":[\"null\",\"headerworks.datatype.UUID\"]}],\"name\":\"Trace\",\"type\":\"record\"}]}],\"name\":\"Data\",\"namespace\":\"headerworks\",\"type\":\"record\"}]"
}

func (_ *UnionHeaderworksData) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionHeaderworksData) SetLong(v int64) {

	r.UnionType = (UnionHeaderworksDataTypeEnum)(v)
}

func (r *UnionHeaderworksData) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		r.HeaderworksData = NewHeaderworksData()
		return &types.Record{Target: (&r.HeaderworksData)}
	}
	panic("Unknown field index")
}
func (_ *UnionHeaderworksData) NullField(i int)                  { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) SetDefault(i int)                 { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *UnionHeaderworksData) Finalize()                        {}

func (r *UnionHeaderworksData) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionHeaderworksDataTypeEnumHeaderworksData:
		return json.Marshal(map[string]interface{}{"headerworks.Data": r.HeaderworksData})
	}
	return nil, fmt.Errorf("invalid value for *UnionHeaderworksData")
}

func (r *UnionHeaderworksData) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["headerworks.Data"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.HeaderworksData)
	}
	return fmt.Errorf("invalid value for *UnionHeaderworksData")
}
