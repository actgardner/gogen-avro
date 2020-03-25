// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     defaults.avsc
 */
package avro

import (
	"io"
	"github.com/actgardner/gogen-avro/vm/types"
	"github.com/actgardner/gogen-avro/vm"
	"github.com/actgardner/gogen-avro/compiler"
)


type UnionRecord struct {

	
	
		A string
	

	
	
		Id *UnionNullInt
	

	
	
		Name *UnionNullString
	

}

const UnionRecordAvroCRC64Fingerprint = "\xfeS\x1bd\xa1\xfc͒"

func NewUnionRecord() (*UnionRecord) {
	return &UnionRecord{}
}

func DeserializeUnionRecord(r io.Reader) (*UnionRecord, error) {
	t := NewUnionRecord()
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

func DeserializeUnionRecordFromSchema(r io.Reader, schema string) (*UnionRecord, error) {
	t := NewUnionRecord()

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

func writeUnionRecord(r *UnionRecord, w io.Writer) error {
	var err error
	
	err = vm.WriteString( r.A, w)
	if err != nil {
		return err
	}
	
	err = writeUnionNullInt( r.Id, w)
	if err != nil {
		return err
	}
	
	err = writeUnionNullString( r.Name, w)
	if err != nil {
		return err
	}
	
	return err
}

func (r *UnionRecord) Serialize(w io.Writer) error {
	return writeUnionRecord(r, w)
}

func (r *UnionRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"a\",\"type\":\"string\"},{\"default\":null,\"name\":\"id\",\"type\":[\"null\",\"int\"]},{\"default\":null,\"name\":\"name\",\"type\":[\"null\",\"string\"]}],\"name\":\"UnionRecord\",\"type\":\"record\"}"
}

func (r *UnionRecord) SchemaName() string {
	return "UnionRecord"
}

func (_ *UnionRecord) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ *UnionRecord) SetInt(v int32) { panic("Unsupported operation") }
func (_ *UnionRecord) SetLong(v int64) { panic("Unsupported operation") }
func (_ *UnionRecord) SetFloat(v float32) { panic("Unsupported operation") }
func (_ *UnionRecord) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionRecord) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ *UnionRecord) SetString(v string) { panic("Unsupported operation") }
func (_ *UnionRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *UnionRecord) Get(i int) types.Field {
	switch (i) {
	
	case 0:
		
		
			return (*types.String)(&r.A)
		
	
	case 1:
		
			r.Id = NewUnionNullInt()

		
		
			return r.Id
		
	
	case 2:
		
			r.Name = NewUnionNullString()

		
		
			return r.Name
		
	
	}
	panic("Unknown field index")
}

func (r *UnionRecord) SetDefault(i int) {
	switch (i) {
	
        
	
        
	case 1:
       	 	r.Id = NewUnionNullInt()

		return
	
	
        
	case 2:
       	 	r.Name = NewUnionNullString()

		return
	
	
	}
	panic("Unknown field index")
}

func (_ *UnionRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UnionRecord) AppendArray() types.Field { panic("Unsupported operation") }
func (_ *UnionRecord) Finalize() { }


func (_ *UnionRecord) AvroCRC64Fingerprint() []byte {
  return []byte(UnionRecordAvroCRC64Fingerprint)
}
