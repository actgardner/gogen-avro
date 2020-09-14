package main

import (
	"reflect"

	"github.com/actgardner/gogen-avro/v7/schema"
)

type compiler struct {
	methods map[string]Method	
}

func Compile(schema []byte, t interface{}) unsafe.Pointer {
	c := compiler{
		methods: make(map[string]Method),
	}
}

func deserializerMethod(t schema.AvroType) uintptr {
	switch v := writer.(type) {
	case *schema.Reference:
		panic("can't JIT references yet")	
	case *schema.MapField:
		panic("can't JIT map yet")	
	case *schema.ArrayField:
		panic("can't JIT array yet")	
	case *schema.UnionField:
		panic("can't JIT union yet")	
	case *schema.IntField:
		
	case *schema.LongField:
		
	case *schema.StringField:
		
	case *schema.BytesField:
		
	case *schema.FloatField:
	
	case *schema.DoubleField:
		
	case *schema.BoolField:
		
	case *schema.NullField:
	}
}

func (c *compiler) allocate(t schema.AvroType) {
	switch v := writer.(type) {
	case *schema.Reference:
			
	case *schema.MapField:
		panic("can't JIT map yet")	
	case *schema.ArrayField:
		panic("can't JIT array yet")	
	case *schema.UnionField:
		panic("can't JIT union yet")
	}
}

func (c *compiler) handleRecord(rd *schema.RecordDefinition, s interface{}) (*program, error) {
	t := reflect.TypeOf(s)
	for i, f := range rd.Fields() {
		offset := t.Field(i).Offset
		switch f.
	}
}
