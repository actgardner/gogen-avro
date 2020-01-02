package templates

import (
	"github.com/actgardner/gogen-avro/generator/namer"
	avro "github.com/actgardner/gogen-avro/schema"
)

func NodeTypeMetadata(n avro.Node) interface{} {
	return n.GetGeneratorMetadata(namer.MetadataKey)
}

func FieldMetadata(f *avro.Field) interface{} {
	return f.GetGeneratorMetadata(namer.MetadataKey)
}
