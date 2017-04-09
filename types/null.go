package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var NullType = &PrimitiveType {
	fieldType: "Null",
	goType: "interface{}",
	serializerMethod: writeNullMethod,
	deserializerMethod: readNullMethod,
	blocks: []generator.Block{},
}

var writeNullMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeNull"},
	Arguments: []*generator.StructField{
		&generator.StructField{"_", "interface{}"},
		&generator.StructField{"_", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
	Body: `
		return nil
	`,
}

var readNullMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readNull"},
	Arguments: []*generator.StructField{
		&generator.StructField{"_", "io.Reader"}
	},
	ReturnTypes: []string{"interface{}", "error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
        Body: ` 
		return nil, nil
	`,
}
`
