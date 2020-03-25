// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     fixed.avsc
 */
package avro

import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)


type FixedTestRecord struct {

	
	
		FixedField TestFixedType
	

	
	
		AnotherFixed TestFixedType
	

}

const FixedTestRecordAvroCRC64Fingerprint = "\xbaK\xf9~\x1f\xc0\xf1R"

func NewFixedTestRecord() (*FixedTestRecord) {
	return &FixedTestRecord{}
}

func DeserializeFixedTestRecord(r io.Reader) (*FixedTestRecord, error) {
	t := NewFixedTestRecord()
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

func DeserializeFixedTestRecordFromSchema(r io.Reader, schema string) (*FixedTestRecord, error) {
	t := NewFixedTestRecord()

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

func writeFixedTestRecord(r *FixedTestRecord, w io.Writer) error {
	var err error
	
	err = writeTestFixedType( r.FixedField, w)
	if err != nil {
		return err
	}
	
	err = writeTestFixedType( r.AnotherFixed, w)
	if err != nil {
		return err
	}
	
	return err
}

func (r *FixedTestRecord) Serialize(w io.Writer) error {
	return writeFixedTestRecord(r, w)
}

func (r *FixedTestRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"FixedField\",\"type\":{\"name\":\"TestFixedType\",\"size\":12,\"type\":\"fixed\"}},{\"name\":\"AnotherFixed\",\"type\":\"TestFixedType\"}],\"name\":\"FixedTestRecord\",\"type\":\"record\"}"
}

func (r *FixedTestRecord) SchemaName() string {
	return "FixedTestRecord"
}

func (_ *FixedTestRecord) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetInt(v int32) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetLong(v int64) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetString(v string) { panic("Unsupported operation") }
func (_ *FixedTestRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *FixedTestRecord) Get(i int) types.Field {
	switch (i) {
	
	case 0:
		
		
			return (*TestFixedTypeWrapper)(&r.FixedField)
		
	
	case 1:
		
		
			return (*TestFixedTypeWrapper)(&r.AnotherFixed)
		
	
	}
	panic("Unknown field index")
}

func (r *FixedTestRecord) SetDefault(i int) {
	switch (i) {
	
        
	
        
	
	}
	panic("Unknown field index")
}

func (_ *FixedTestRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *FixedTestRecord) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *FixedTestRecord) Finalize() { }


func (_ *FixedTestRecord) AvroCRC64Fingerprint() []byte {
  return []byte(FixedTestRecordAvroCRC64Fingerprint)
}
