// Wrappers for Avro primitive types implementing the methods required by GADGT
package types

// The interface neeed by GADGT to enter and set fields on a type
// Most types only need to implement a subset
type Field interface {
	// Assign a primitive field
	SetBoolean(v bool)
	SetInt(v int32)
	SetLong(v int64)
	SetFloat(v float32)
	SetDouble(v float64)
	SetBytes(v []byte)
	SetString(v string)

	// Get a nested field
	Get(i int) Field
	// Set the default value for a given field
	SetDefault(i int)

	// Append a new value to a map or array and enter it
	AppendMap(key string) Field
	AppendArray() Field

	// Finalize a field if necessary
	Finalize()
}
