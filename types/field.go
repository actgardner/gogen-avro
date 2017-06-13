package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

type Field struct {
	avroName   string
	avroType   AvroType
	defValue   interface{}
	hasDef     bool
	definition map[string]interface{}
}

func NewField(avroName string, avroType AvroType, defValue interface{}, hasDef bool, definition map[string]interface{}) *Field {
	return &Field{
		avroName:   avroName,
		avroType:   avroType,
		defValue:   defValue,
		hasDef:     hasDef,
		definition: definition,
	}
}

func (f *Field) Name() string {
	return f.avroName
}

func (f *Field) GoName() string {
	return generator.ToPublicName(f.avroName)
}

func (f *Field) HasDefault() bool {
	return f.defValue == nil
}

func (f *Field) Default() interface{} {
	return f.defValue
}

func (f *Field) Type() AvroType {
	return f.avroType
}

func (f *Field) Definition(scope map[QualifiedName]interface{}) (map[string]interface{}, error) {
	var err error
	f.definition["type"], err = f.avroType.Definition(scope)
	if err != nil {
		return nil, err
	}

	return f.definition, nil
}
