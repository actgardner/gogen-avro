package main

import (
	"reflect"
	"unsafe"

	"github.com/actgardner/gogen-avro/v7/schema"
)

type compiler struct {
	methods map[string]Method
}

func Compile(schema []byte, t interface{}) unsafe.Pointer {
	c := compiler{
		methods: make(map[string]Method),
	}
	return nil
}

func deserializerMethod(t schema.AvroType) (uintptr, error) {
	switch v := t.(type) {
	case *schema.Reference:
		panic("can't JIT references yet")
	case *schema.MapField:
		panic("can't JIT map yet")
	case *schema.ArrayField:
		panic("can't JIT array yet")
	case *schema.UnionField:
		panic("can't JIT union yet")
	case *schema.IntField:
		x := assignIntToInt
		return uintptr(unsafe.Pointer(&x)), nil
	case *schema.LongField:
		x := assignLongToLong
		return uintptr(unsafe.Pointer(&x)), nil
	case *schema.StringField:
		x := assignStringToString
		return uintptr(unsafe.Pointer(&x)), nil
	case *schema.BytesField:
		x := assignBytesToBytes
		return uintptr(unsafe.Pointer(&x)), nil
	case *schema.FloatField:
		x := assignFloatToFloat
		return uintptr(unsafe.Pointer(&x)), nil
	case *schema.DoubleField:
		x := assignDoubleToDouble
		return uintptr(unsafe.Pointer(&x)), nil
	case *schema.BoolField:
		x := assignBoolToBool
		return uintptr(unsafe.Pointer(&x)), nil
	}
	return 0, nil
}

func (c *compiler) allocate(t schema.AvroType) (uintptr, error) {
	switch v := t.(type) {
	case *schema.Reference:
		switch ref := v.Def.(type) {
		case *schema.RecordDefinition:
			return c.handleRecord(ref).FuncPtr(), nil
		}
	case *schema.MapField:
		panic("can't JIT map yet")
	case *schema.ArrayField:
		panic("can't JIT array yet")
	case *schema.UnionField:
		panic("can't JIT union yet")
	}
}

func (c *compiler) handleRecord(rd *schema.RecordDefinition, s interface{}) (*Method, error) {
	m, err := NewMethod()
	if err != nil {
		return nil, err
	}
	c.methods[rd.GoType] = m

	t := reflect.TypeOf(s)
	for i, f := range rd.Fields() {
		offset := t.Field(i).Offset
		deserializer, err := c.deserializerMethod(f.Type())
		if err != nil {
			return nil, err
		}
		// Deserializer signature is (reader, uintptr) -> (error), allocate 48 bytes
		m.AppendOp(SubqImm(Rax, 0x30))
		m.AppendOp(MovqDestIDRSP(Rbp, 0x28))
		m.AppendOp(LeaqSrcIDRSP(Rbp, 0x28))
		m.AppendOp(MovqSrcIDRSP(Rax, 0x0))
		m.AppendOp(MovqDestIDRSP(Rax, 0x0))
		m.AppendCall(deserializer)
	}
}
