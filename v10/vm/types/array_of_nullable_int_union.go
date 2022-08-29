package types

type ArrayOfNullableIntUnion struct {
	Target *[]*int32
}

func (b ArrayOfNullableIntUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableIntUnion) AppendArray() Field {
	var v int32

	*b.Target = append(*b.Target, &v)

	return &Int{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableIntUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableIntUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*int32, 0, s)
	}
}

func (b ArrayOfNullableIntUnion) Finalize() {}
