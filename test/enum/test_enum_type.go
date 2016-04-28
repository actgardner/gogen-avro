package avro

type TestEnumType int32

const (
	TestSymbol1 TestEnumType = 0
	TestSymbol2 TestEnumType = 1
	TestSymbol3 TestEnumType = 2
)

func (e TestEnumType) String() string {
	switch e {
	case TestSymbol1:
		return "TestSymbol1"
	case TestSymbol2:
		return "testSymbol2"
	case TestSymbol3:
		return "testSymbol3"

	}
	return "Unknown"
}
