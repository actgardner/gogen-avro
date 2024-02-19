package types

type ByteArray []byte

type ArrayOfNullableBytesUnion struct {
	Target *[]*ByteArray
}

func (b ArrayOfNullableBytesUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableBytesUnion field")
}

func (b ArrayOfNullableBytesUnion) AppendArray() Field {
	var v ByteArray

	*b.Target = append(*b.Target, &v)

	return &Bytes{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableBytesUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableBytesUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*ByteArray, 0, s)
	}
}

func (b ArrayOfNullableBytesUnion) Finalize() {}
