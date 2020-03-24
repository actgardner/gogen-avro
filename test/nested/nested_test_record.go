// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     nested.avsc
 */
package avro

import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)


type NestedTestRecord struct {

	
	
		NumberField *NumberRecord
	

	
	
		OtherField *NestedRecord
	

}

var NestedTestRecordAvroCRC64Fingerprint = []byte{0x62,0x7b,0x6d,0x5c,0x44,0xbe,0xaa,0x96}

func NewNestedTestRecord() (*NestedTestRecord) {
	return &NestedTestRecord{}
}

func DeserializeNestedTestRecord(r io.Reader) (*NestedTestRecord, error) {
	t := NewNestedTestRecord()
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

func DeserializeNestedTestRecordFromSchema(r io.Reader, schema string) (*NestedTestRecord, error) {
	t := NewNestedTestRecord()

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

func writeNestedTestRecord(r *NestedTestRecord, w io.Writer) error {
	var err error
	
	err = writeNumberRecord( r.NumberField, w)
	if err != nil {
		return err
	}
	
	err = writeNestedRecord( r.OtherField, w)
	if err != nil {
		return err
	}
	
	return err
}

func (r *NestedTestRecord) Serialize(w io.Writer) error {
	return writeNestedTestRecord(r, w)
}

func (r *NestedTestRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"NumberField\",\"type\":{\"fields\":[{\"name\":\"IntField\",\"type\":\"int\"},{\"name\":\"LongField\",\"type\":\"long\"},{\"name\":\"FloatField\",\"type\":\"float\"},{\"name\":\"DoubleField\",\"type\":\"double\"}],\"name\":\"NumberRecord\",\"type\":\"record\"}},{\"name\":\"OtherField\",\"type\":{\"fields\":[{\"name\":\"StringField\",\"type\":\"string\"},{\"name\":\"BoolField\",\"type\":\"boolean\"},{\"name\":\"BytesField\",\"type\":\"bytes\"}],\"name\":\"NestedRecord\",\"type\":\"record\"}}],\"name\":\"NestedTestRecord\",\"type\":\"record\"}"
}

func (r *NestedTestRecord) SchemaName() string {
	return "NestedTestRecord"
}

func (_ *NestedTestRecord) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetInt(v int32) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetLong(v int64) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetString(v string) { panic("Unsupported operation") }
func (_ *NestedTestRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *NestedTestRecord) Get(i int) types.Field {
	switch (i) {
	
	case 0:
		
			r.NumberField = NewNumberRecord()

		
		
			return r.NumberField
		
	
	case 1:
		
			r.OtherField = NewNestedRecord()

		
		
			return r.OtherField
		
	
	}
	panic("Unknown field index")
}

func (r *NestedTestRecord) SetDefault(i int) {
	switch (i) {
	
        
	
        
	
	}
	panic("Unknown field index")
}

func (_ *NestedTestRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *NestedTestRecord) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *NestedTestRecord) Finalize() { }


func (_ *NestedTestRecord) AvroCRC64Fingerprint() []byte {
  return NestedTestRecordAvroCRC64Fingerprint
}
