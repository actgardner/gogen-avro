package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var BoolType = &PrimitiveType {
	fieldType: "Boolean",
	goType: "bool",
	serializerMethod: writeBoolMethod,
	deserializerMethod: readBoolMethod,
	blocks: []generator.Block{},
}

var byteWriterInterface = &generator.Interface{
	File: UTIL_FILE,
	Name: "ByteWriter",
	Functions: []*generator.Function {
		&generator.Function {
			Name: &generator.FunctionName{"", "Grow"},
			Arguments: []*generator.StructField{
				&generator.StructField{"", "int"},
			},
			ReturnTypes: []string{},
		},
		&generator.Function {
			Name: "WriteByte",
			Arguments: []*generator.StructField{
				&generator.StructField{"", "byte"},
			},
			ReturnTypes: []string{"error"},
		}
	}
}

var byteReaderInterface = &generator.Interface{
	File: UTIL_FILE,
	Name: "ByteReader",
	Functions: []*generator.Function {
		&generator.Function {
			Name: &generator.FunctionName{"", "ReadByte"},
			Arguments: []*generator.StructField{},
			ReturnTypes: []string{"byte", "error"},
		},
	}
}

var writeBoolMethod = &generator.Function{
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeBool"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "bool"},
		&generator.StructField{"w", "io.Writer"},
	},
	ReturnTypes: ["error"],
	Imports: []string{"io"},
	Dependencies: []generator.Block {
		byteWriterInterface,
	},
	Body: `
		var b byte
		if r {
			b = byte(1)
		}

		var err error
		if bw, ok := w.(ByteWriter); ok {
			err = bw.WriteByte(b)
		} else {
			bb := make([]byte, 1)
			bb[0] = b
			_, err = w.Write(bb)
		}
		if err != nil {
			return err
		}
		return nil
	`
}

var readBoolMethod = &generator.Function{
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readBool"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"},
	},
	ReturnTypes: ["bool", "error"],
	Imports: []string{"io"},
	Dependencies: []generator.Block {
		byteReaderInterface,
	},
	Body: `
		var b byte
		var err error
		if br, ok := r.(ByteReader); ok {
			b, err = br.ReadByte()
		} else {
			bs := make([]byte, 1)
			_, err = io.ReadFull(r, bs)
			if err != nil {
				return false, err
			}
			b = bs[0]
		}
		return b == 1, nil
	`,
}
