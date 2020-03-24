// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     stringtest.avsc
 */
package avro

import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)


type StringRec struct {

	
	
		ProductName string
	

}

var StringRecAvroCRC64Fingerprint = []byte{0x77,0x83,0x36,0xab,0x9d,0xe9,0x0,0x15}

func NewStringRec() (*StringRec) {
	return &StringRec{}
}

func DeserializeStringRec(r io.Reader) (*StringRec, error) {
	t := NewStringRec()
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

func DeserializeStringRecFromSchema(r io.Reader, schema string) (*StringRec, error) {
	t := NewStringRec()

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

func writeStringRec(r *StringRec, w io.Writer) error {
	var err error
	
	err = vm.WriteString( r.ProductName, w)
	if err != nil {
		return err
	}
	
	return err
}

func (r *StringRec) Serialize(w io.Writer) error {
	return writeStringRec(r, w)
}

func (r *StringRec) Schema() string {
	return "{\"fields\":[{\"name\":\"productName\",\"type\":\"string\"}],\"name\":\"StringRec\",\"type\":\"record\"}"
}

func (r *StringRec) SchemaName() string {
	return "StringRec"
}

func (_ *StringRec) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *StringRec) SetInt(v int32) { panic("Unsupported operation") }
func (_ *StringRec) SetLong(v int64) { panic("Unsupported operation") }
func (_ *StringRec) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *StringRec) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *StringRec) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *StringRec) SetString(v string) { panic("Unsupported operation") }
func (_ *StringRec) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *StringRec) Get(i int) types.Field {
	switch (i) {
	
	case 0:
		
		
			return (*types.String)(&r.ProductName)
		
	
	}
	panic("Unknown field index")
}

func (r *StringRec) SetDefault(i int) {
	switch (i) {
	
        
	
	}
	panic("Unknown field index")
}

func (_ *StringRec) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *StringRec) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *StringRec) Finalize() { }


func (_ *StringRec) AvroCRC64Fingerprint() []byte {
  return StringRecAvroCRC64Fingerprint
}
