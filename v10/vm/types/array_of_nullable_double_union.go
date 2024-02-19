package types

type ArrayOfNullableDoubleUnion struct {
	Target *[]*float64
}

func (b ArrayOfNullableDoubleUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableDoubleUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableDoubleUnion field")
}

func (b ArrayOfNullableDoubleUnion) AppendArray() Field {
	var v float64

	*b.Target = append(*b.Target, &v)

	return &Double{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableDoubleUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableDoubleUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*float64, 0, s)
	}
}

func (b ArrayOfNullableDoubleUnion) Finalize() {}
