package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var FloatType = &PrimitiveType {
	fieldType: "Float",
	goType: "float32",
	serializerMethod: writeFloatMethod,
	deserializerMethod: readFloatMethod,
	blocks: []generator.Block{},
}

var writeFloatMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeFloat"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "float32"},
		&generator.StructField{"w", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{
		encodeFloatMethod,
	},
        Body: `
		bits := uint64(math.Float32bits(r))
		const byteCount = 4
		return encodeFloat(w, byteCount, bits)
	`,
}

var encodeFloatMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "encodeFloat"},
	Arguments: []*generator.StructField{
		&generator.StructField{"w", "io.Writer"}
		&generator.StructField{"byteCount", "int"}
		&generator.StructField{"bits", "uint64"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{
		byteWriterInterface,
	},
        Body: `
		var err error
		var bb []byte
		bw, ok := w.(ByteWriter)
		if ok {
			bw.Grow(byteCount)
		} else {
			bb = make([]byte, 0, byteCount)
		}
		for i := 0; i < byteCount; i++ {
			if bw != nil {
				err = bw.WriteByte(byte(bits & 255))
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, byte(bits&255))
			}
			bits = bits >> 8
		}
		if bw == nil {
			_, err = w.Write(bb)
			return err
		}
		return nil
	`,
}

var readFloatMethod = &generator.Function {
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readFloat"},
	Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"}
	},
	ReturnTypes: []string{"float32", "error"},
	Imports: []string{"io", "encoding/binary", "math"},
	Dependencies: []generator.Block{},
        Body: `
		buf := make([]byte, 4)
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}
		bits := binary.LittleEndian.Uint32(buf)
		val := math.Float32frombits(bits)
		return val, nil
	`,
}
