package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var DoubleType = &PrimitiveType {
	fieldType: "Double",
	goType: "float64",
	serializerMethod: writeDoubleMethod,
	deserializerMethod: readDoubleMethod,
	blocks: []generator.Block{},
}

var writeDoubleMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeDouble"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "float64"},
		&generator.StructField{"w", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io", "math"},
	Dependencies: []generator.Block{
		encodeFloatMethod,
	},
        Body: `
		bits := uint64(math.Float64bits(r))
		const byteCount = 8
		return encodeFloat(w, byteCount, bits)
	`,
}

var readDoubleMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readDouble"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"}
	},
	ReturnTypes: []string{"float64", "error"},
	Imports: []string{"io", "encoding/binary", "math"},
	Dependencies: []generator.Block{},
        Body: ` 
		buf := make([]byte, 8)
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}
		bits := binary.LittleEndian.Uint64(buf)
		val := math.Float64frombits(bits)
		return val, nil
	`,
}
