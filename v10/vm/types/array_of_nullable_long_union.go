package types

type ArrayOfNullableLongUnion struct {
	Target *[]*int64
}

func (b ArrayOfNullableLongUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableLongUnion field")
}

func (b ArrayOfNullableLongUnion) AppendArray() Field {
	var v int64

	*b.Target = append(*b.Target, &v)

	return &Long{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableLongUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableLongUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*int64, 0, s)
	}
}

func (b ArrayOfNullableLongUnion) Finalize() {}
