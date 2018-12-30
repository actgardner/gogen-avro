package types

type Field interface {
	// Assign a primitive field
	SetBoolean(v bool)
	SetInt(v int32)
	SetLong(v int64)
	SetFloat(v float32)
	SetDouble(v float64)
	SetBytes(v []byte)
	SetString(v string)

	// Get a nested field definition
	Get(i int) Field

	// Append a new value to a map or array
	AppendMap(key string) Field
	AppendArray() Field
}
