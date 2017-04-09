package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

var IntType = &PrimitiveType {
	fieldType: "Int",
	goType: "int32",
	serializerMethod: writeIntMethod,
	deserializerMethod: readIntMethod,
	blocks: []generator.Block{},
}

var writeIntMethod = &generator.Function{
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "writeInt"},
        Arguments: []*generator.StructField{
		&generator.StructField{"r", "int32"},
		&generator.StructField{"w", "io.Writer"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{encodeIntMethod},
        Body: `
		downShift := uint32(31)
		encoded := uint64((uint32(r) << 1) ^ uint32(r >> downShift))
		const maxByteSize = 5
		return encodeInt(w, maxByteSize, encoded)
	`,
}

var encodeIntMethod = &generator.Function{
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "encodeInt"},
        Arguments: []*generator.StructField{
		&generator.StructField{"w", "io.Writer"}, 
		&generator.StructField{"byteCount", "int"}, 
		&generator.StructField{"encoded", "uint64"}
	},
	ReturnTypes: []string{"error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{byteWriterInterface},
	Body: `
		var err error
		var bb []byte
		bw, ok := w.(ByteWriter)
		if ok {
			bw.Grow(byteCount)
		} else {
			bb = make([]byte, 0, byteCount)
		}

		if encoded == 0 {
			if bw != nil {
				err = bw.WriteByte(0)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, byte(0))
			}
		} else {
			for encoded > 0 {
				b := byte(encoded & 127)
				encoded = encoded >> 7
				if !(encoded == 0) {
					b |= 128
				}
				if bw != nil {
					err = bw.WriteByte(b)
					if err != nil {
						return err
					}
				} else {
					bb = append(bb, b)
				}
			}
		}
		if bw == nil {
			_, err := w.Write(bb)
			return err
		}
		return nil
	`
}

var readIntMethod = &generator.Function{
	File: UTIL_FILE,
	Name: &generator.FunctionName{"", "readInt"},
        Arguments: []*generator.StructField{
		&generator.StructField{"r", "io.Reader"}
	},
	ReturnTypes: []string{"int32", "error"},
	Imports: []string{"io"},
	Dependencies: []generator.Block{},
	Body: `
		var v int
		buf := make([]byte, 1)
		for shift := uint(0); ; shift += 7 {
			if _, err := io.ReadFull(r, buf); err != nil {
				return 0, err
			}
			b := buf[0]
			v |= int(b&127) << shift
			if b&128 == 0 {
				break
			}
		}
		datum := (int32(v>>1) ^ -int32(v&1))
		return datum, nil
	`
}
