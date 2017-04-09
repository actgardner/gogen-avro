package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var LongType = &PrimitiveType {
	fieldType: "Lonng",
	goType: "int64",
	serializerMethod: writeLongMethod,
	deserializerMethod: readLongMethod,
	blocks: []generator.Block{},
}

var writeLongMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeLong"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "int64"},
		&generator.StructField{"w", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{
		encodeIntMethod,
	},
        Body: `
		downShift := uint64(63)
		encoded := uint64((r << 1) ^ (r >> downShift))
		const maxByteSize = 10
		return encodeInt(w, maxByteSize, encoded)
	`,
}

var readDoubleMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readDouble"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"}
	},
	ReturnTypes: []string{"int64", "error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
        Body: `
		var v uint64
		buf := make([]byte, 1)
		for shift := uint(0); ; shift += 7 {
			if _, err := io.ReadFull(r, buf); err != nil {
				return 0, err
			}
			b := buf[0]
			v |= uint64(b&127) << shift
			if b&128 == 0 {
				break
			}
		}
		datum := (int64(v>>1) ^ -int64(v&1))
		return datum, nil
	`
}
