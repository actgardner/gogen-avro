package main

import (
	"io"
)

type Assigner func(r io.Reader, f types.Field) error

func assign(r io.Reader, record *primitive.PrimitiveTestRecord) error {
	err := assignIntToInt(r, &f.IntField)
	if err != nil {
		return err
	}

	err = assignLongToLong(r, &f.LongField)
	if err != nil {
		return err
	}

	err = assignFloatToFloat(r, &f.FloatField)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	p, err := newProgram()
	if err != nil {
		fmt.Printf("mmap error: %v", err)
	}
	p.callAssigner(reflect.ValueOf(assignBool).Pointer())

	//fn := p.funcPtr()
	//exeFn := *(*Assigner)(fn)

	src := primitive.PrimitiveTestRecord{
		IntField:    1,
		LongField:   2,
		FloatField:  3.2,
		DoubleField: 4.1,
		BoolField:   true,
		StringField: "string",
	}

	var buf bytes.Buffer
	err = src.Serialize(&buf)
	//var target primitive.PrimitiveTestRecord

	//err = exeFn(buf, &target)
	var target primitive.PrimitiveTestRecord
	f := assign(&buf, &target)
	fmt.Printf("Result: %v\n", f)
}
