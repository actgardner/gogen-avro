package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var BytesType = &PrimitiveType {
	fieldType: "Bytes",
	goType: "[]byte",
	serializerMethod: writeBytesMethod,
	deserializerMethod: readBytesMethod,
	blocks: []generator.Block{},
}

var writeBytesMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeBytes"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "[]byte"}
		&generator.StructField{"w", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
	Body: `
		err := writeLong(int64(len(r)), w)
		if err != nil {
			return err
		}
		_, err = w.Write(r)
		return err
	`,
}

var readBytesMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readBytes"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"}
	},
	ReturnTypes: []string{"[]byte", "error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
	Body: `
		size, err := readLong(r)
		if err != nil {
			return nil, err
		}
		bb := make([]byte, size)
		_, err = io.ReadFull(r, bb)
		return bb, err
	`,
}
