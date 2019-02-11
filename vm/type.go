package vm

type Type int

const (
	Unused Type = iota
	Null
	Boolean
	Int
	Long
	Float
	Double
	Bytes
	String
	MapKey
	UnionElem
	Fixed
)

func (t Type) String() string {
	switch t {
	case Null:
		return "Null"
	case Boolean:
		return "Boolean"
	case Int:
		return "Int"
	case Long:
		return "Long"
	case Float:
		return "Float"
	case Double:
		return "Double"
	case Bytes:
		return "Bytes"
	case String:
		return "String"
	case MapKey:
		return "MapKey"
	case UnionElem:
		return "UnionElem"
	case Fixed:
		return "Fixed"
	}
	return "-"
}
