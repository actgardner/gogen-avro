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

var _ = fmt.Printf

type PrimitiveTestRecord struct {
	IntField int32 `json:"IntField"`

	LongField int64 `json:"LongField"`

	FloatField float32 `json:"FloatField"`

	DoubleField float64 `json:"DoubleField"`

	StringField string `json:"StringField"`

	BoolField bool `json:"BoolField"`

	BytesField Bytes `json:"BytesField"`
}

const PrimitiveTestRecordAvroCRC64Fingerprint = "۲\x16\xe9\xfet@\x13"

func NewPrimitiveTestRecord() PrimitiveTestRecord {
	r := PrimitiveTestRecord{}
	r.IntField = 1.2345689e+07
	r.LongField = 2.3456789e+08
	r.FloatField = 1e+08
	r.DoubleField = 800000
	r.StringField = "defaultstring"
	r.BoolField = true
	r.BytesField = []byte("\x04\x01\x05\xfd")
	return r
}

func DeserializePrimitiveTestRecord(r io.Reader) (PrimitiveTestRecord, error) {
	t := NewPrimitiveTestRecord()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializePrimitiveTestRecordFromSchema(r io.Reader, schema string) (PrimitiveTestRecord, error) {
	t := NewPrimitiveTestRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writePrimitiveTestRecord(r PrimitiveTestRecord, w io.Writer) error {
	var err error
	err = vm.WriteInt(r.IntField, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.LongField, w)
	if err != nil {
		return err
	}
	err = vm.WriteFloat(r.FloatField, w)
	if err != nil {
		return err
	}
	err = vm.WriteDouble(r.DoubleField, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.StringField, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.BoolField, w)
	if err != nil {
		return err
	}
	err = vm.WriteBytes(r.BytesField, w)
	if err != nil {
		return err
	}
	return err
}

func (r PrimitiveTestRecord) Serialize(w io.Writer) error {
	return writePrimitiveTestRecord(r, w)
}

func (r PrimitiveTestRecord) Schema() string {
	return "{\"fields\":[{\"default\":12345689,\"name\":\"IntField\",\"type\":\"int\"},{\"default\":234567890,\"name\":\"LongField\",\"type\":\"long\"},{\"default\":100000000,\"name\":\"FloatField\",\"type\":\"float\"},{\"default\":800000,\"name\":\"DoubleField\",\"type\":\"double\"},{\"default\":\"defaultstring\",\"name\":\"StringField\",\"type\":\"string\"},{\"default\":true,\"name\":\"BoolField\",\"type\":\"boolean\"},{\"default\":\"\\u0004\\u0001\\u0005ý\",\"name\":\"BytesField\",\"type\":\"bytes\"}],\"name\":\"PrimitiveTestRecord\",\"type\":\"record\"}"
}

func (r PrimitiveTestRecord) SchemaName() string {
	return "PrimitiveTestRecord"
}

func (_ PrimitiveTestRecord) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetInt(v int32)       { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetLong(v int64)      { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetString(v string)   { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *PrimitiveTestRecord) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.Int{Target: &r.IntField}

		return w

	case 1:
		w := types.Long{Target: &r.LongField}

		return w

	case 2:
		w := types.Float{Target: &r.FloatField}

		return w

	case 3:
		w := types.Double{Target: &r.DoubleField}

		return w

	case 4:
		w := types.String{Target: &r.StringField}

		return w

	case 5:
		w := types.Boolean{Target: &r.BoolField}

		return w

	case 6:
		w := BytesWrapper{Target: &r.BytesField}

		return w

	}
	panic("Unknown field index")
}

func (r *PrimitiveTestRecord) SetDefault(i int) {
	switch i {
	case 0:
		r.IntField = 1.2345689e+07
		return
	case 1:
		r.LongField = 2.3456789e+08
		return
	case 2:
		r.FloatField = 1e+08
		return
	case 3:
		r.DoubleField = 800000
		return
	case 4:
		r.StringField = "defaultstring"
		return
	case 5:
		r.BoolField = true
		return
	case 6:
		r.BytesField = []byte("\x04\x01\x05\xfd")
		return
	}
	panic("Unknown field index")
}

func (r *PrimitiveTestRecord) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ PrimitiveTestRecord) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) HintSize(int)                     { panic("Unsupported operation") }
func (_ PrimitiveTestRecord) Finalize()                        {}

func (_ PrimitiveTestRecord) AvroCRC64Fingerprint() []byte {
	return []byte(PrimitiveTestRecordAvroCRC64Fingerprint)
}

func (r PrimitiveTestRecord) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["IntField"], err = json.Marshal(r.IntField)
	if err != nil {
		return nil, err
	}
	output["LongField"], err = json.Marshal(r.LongField)
	if err != nil {
		return nil, err
	}
	output["FloatField"], err = json.Marshal(r.FloatField)
	if err != nil {
		return nil, err
	}
	output["DoubleField"], err = json.Marshal(r.DoubleField)
	if err != nil {
		return nil, err
	}
	output["StringField"], err = json.Marshal(r.StringField)
	if err != nil {
		return nil, err
	}
	output["BoolField"], err = json.Marshal(r.BoolField)
	if err != nil {
		return nil, err
	}
	output["BytesField"], err = json.Marshal(r.BytesField)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *PrimitiveTestRecord) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["IntField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.IntField); err != nil {
			return err
		}
	} else {
		r.IntField = 1.2345689e+07
	}
	val = func() json.RawMessage {
		if v, ok := fields["LongField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.LongField); err != nil {
			return err
		}
	} else {
		r.LongField = 2.3456789e+08
	}
	val = func() json.RawMessage {
		if v, ok := fields["FloatField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.FloatField); err != nil {
			return err
		}
	} else {
		r.FloatField = 1e+08
	}
	val = func() json.RawMessage {
		if v, ok := fields["DoubleField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.DoubleField); err != nil {
			return err
		}
	} else {
		r.DoubleField = 800000
	}
	val = func() json.RawMessage {
		if v, ok := fields["StringField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.StringField); err != nil {
			return err
		}
	} else {
		r.StringField = "defaultstring"
	}
	val = func() json.RawMessage {
		if v, ok := fields["BoolField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.BoolField); err != nil {
			return err
		}
	} else {
		r.BoolField = true
	}
	val = func() json.RawMessage {
		if v, ok := fields["BytesField"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.BytesField); err != nil {
			return err
		}
	} else {
		r.BytesField = []byte("\x04\x01\x05\xfd")
	}
	return nil
}