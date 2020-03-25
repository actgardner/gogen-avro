// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     namespace.avsc
 */
package avro

import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)

// A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014
type UUID struct {

	
	
		Uuid string
	

}

const UUIDAvroCRC64Fingerprint = "\xfc\xa43\x98\xee\xe0p\xe2"

func NewUUID() (*UUID) {
	return &UUID{}
}

func DeserializeUUID(r io.Reader) (*UUID, error) {
	t := NewUUID()
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

func DeserializeUUIDFromSchema(r io.Reader, schema string) (*UUID, error) {
	t := NewUUID()

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

func writeUUID(r *UUID, w io.Writer) error {
	var err error
	
	err = vm.WriteString( r.Uuid, w)
	if err != nil {
		return err
	}
	
	return err
}

func (r *UUID) Serialize(w io.Writer) error {
	return writeUUID(r, w)
}

func (r *UUID) Schema() string {
	return "{\"doc\":\"A Universally Unique Identifier, in canonical form in lowercase. Example: de305d54-75b4-431b-adb2-eb6b9e546014\",\"fields\":[{\"default\":\"\",\"name\":\"uuid\",\"type\":\"string\"}],\"name\":\"bodyworks.datatype.UUID\",\"type\":\"record\"}"
}

func (r *UUID) SchemaName() string {
	return "bodyworks.datatype.UUID"
}

func (_ *UUID) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *UUID) SetInt(v int32) { panic("Unsupported operation") }
func (_ *UUID) SetLong(v int64) { panic("Unsupported operation") }
func (_ *UUID) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *UUID) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UUID) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *UUID) SetString(v string) { panic("Unsupported operation") }
func (_ *UUID) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *UUID) Get(i int) types.Field {
	switch (i) {
	
	case 0:
		
		
			return (*types.String)(&r.Uuid)
		
	
	}
	panic("Unknown field index")
}

func (r *UUID) SetDefault(i int) {
	switch (i) {
	
        
	case 0:
       	 	r.Uuid = ""
		return
	
	
	}
	panic("Unknown field index")
}

func (_ *UUID) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UUID) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *UUID) Finalize() { }


func (_ *UUID) AvroCRC64Fingerprint() []byte {
  return []byte(UUIDAvroCRC64Fingerprint)
}
