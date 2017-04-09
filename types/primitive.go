package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

type PrimitiveType struct {
	fieldType                    string
	goType                       string
	serializer             generator.Function
	deserializer           generator.Function

	blocks []generator.Block
}

func (p *PrimitiveType) FieldType() string {
	return p.fieldType
}

func (p *PrimitiveType) GoType() string {
	return p.goType
}

func (p *PrimitiveType) Serializer() generator.Function {
	return p.serializer
}

func (p *PrimitiveType) Deserializer() generator.Function {
	return p.deserializer
}

func (p *PrimitiveType) AddToPkg(pkg *generator.Package) {
	pkg.AddBlocks(p.serializer, p.deserializer, p.blocks...)
}


