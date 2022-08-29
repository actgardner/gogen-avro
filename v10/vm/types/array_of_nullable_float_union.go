package types

type ArrayOfNullableFloatUnion struct {
	Target *[]*float32
}

func (b ArrayOfNullableFloatUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableFloatUnion field")
}

func (b ArrayOfNullableFloatUnion) AppendArray() Field {
	var v float32

	*b.Target = append(*b.Target, &v)

	return &Float{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableFloatUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableFloatUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*float32, 0, s)
	}
}

func (b ArrayOfNullableFloatUnion) Finalize() {}
