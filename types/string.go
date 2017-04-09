package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var StringType = &PrimitiveType {
	fieldType: "String",
	goType: "string",
	serializerMethod: writeStringMethod,
	deserializerMethod: readStringMethod,
	blocks: []generator.Block{},
}

var stringWriterInterface = &generator.Interface{
	File: UTIL_FILE,
	Name: "StringWriter",
	Functions: []*Function{
		&generator.Function{
			File: UTIL_FILE,
			Name: "WriteString",
			Arguments: []*generator.StructField {
				&generator.StructField{"", "string"},
			},
			ReturnTypes: []string{"int", "error"},
			Imports: []string{},
			Dependencies: []generator.Block{},
			Body: "",
		},
	},
}

var writeStringMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeString"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "string"},
		&generator.StructField{"w", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{
		writeLongMethod,
		stringWriterInterface,
	},
        Body: `
		err := writeLong(int64(len(r)), w)
		if err != nil {
			return err
		}
		if sw, ok := w.(StringWriter); ok {
			_, err = sw.WriteString(r)
		} else {
			_, err = w.Write([]byte(r))
		}
		return err
	`,
}

var readStringMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readString"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"}
	},
	ReturnTypes: []string{"string", "error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
        Body: `
		len, err := readLong(r)
		if err != nil {
			return "", err
		}
		bb := make([]byte, len)
		_, err = io.ReadFull(r, bb)
		if err != nil {
			return "", err
		}
		return string(bb), nil
	`
}
