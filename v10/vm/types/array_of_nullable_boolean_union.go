package types

type ArrayOfNullableBooleanUnion struct {
	Target *[]*bool
}

func (b ArrayOfNullableBooleanUnion) SetBoolean(v bool) {
	panic("Unable to assign boolean to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetInt(v int32) {
	panic("Unable to assign boolean to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetLong(v int64) {
	panic("Unable to assign long to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetFloat(v float32) {
	panic("Unable to assign float to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetUnionElem(v int64) {
	panic("Unable to assign union elem to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetDouble(v float64) {
	panic("Unable to assign double to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetBytes(v []byte) {
	panic("Unable to assign bytes to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetString(v string) {
	panic("Unable to assign string to ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) Get(i int) Field {
	panic("Unable to get field from ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) SetDefault(i int) {
	panic("Unable to set default on ArrayOfNullableBooleanUnion field")
}

func (b ArrayOfNullableBooleanUnion) AppendMap(key string) Field {
	panic("Unable to append map key to from ArrayOfNullableIntUnion field")
}

func (b ArrayOfNullableBooleanUnion) AppendArray() Field {
	var v bool

	*b.Target = append(*b.Target, &v)

	return &Boolean{Target: (*b.Target)[len(*b.Target)-1]}
}

func (b ArrayOfNullableBooleanUnion) NullField(i int) {
	(*b.Target)[len(*b.Target)-1] = nil
}

func (b ArrayOfNullableBooleanUnion) HintSize(s int) {
	if len(*b.Target) == 0 {
		*b.Target = make([]*bool, 0, s)
	}
}

func (b ArrayOfNullableBooleanUnion) Finalize() {}
