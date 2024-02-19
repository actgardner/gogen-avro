package types

type ArrayOfNullableStringUnion struct {
	Target *[]*string
}

func (b ArrayOfNullableStringUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableStringUnion field")
}

func (b ArrayOfNullableStringUnion) AppendArray() Field {
	var v string

	*b.Target = append(*b.Target, &v)

	return &String{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableStringUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableStringUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*string, 0, s)
	}
}

func (b ArrayOfNullableStringUnion) Finalize() {}
