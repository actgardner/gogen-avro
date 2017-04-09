package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

type Field struct {
	Name string
	Default interface{}
	HasDefault bool
	Type *AvroType
}

func (f *Field) AvroName() string {
	return f.name
}

func (f *Field) GoName() string {
	return generator.ToPublicName(f.Name)
}

func (f *Field) ResolveReferences(ns *Namespace) error {
	return nil
}

func (f *Field) Schema(names map[QualifiedName]interface{}) interface{} {
	return nil
}
