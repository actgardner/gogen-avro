package vm

import (
	"fmt"
	"io"
)

type Frame struct {
	Target Assignable

	Boolean bool
	Int     int32
	Long    int64
	Float   float32
	Double  float64
	Bytes   []byte
	String  string

	MapKey     string
	Length     int64
	UnionElem  int64
	BlockStart int
}

func Eval(r io.Reader, program []Instruction, target Assignable) (err error) {
	stack := make([]Frame, 256)
	stack[0].Target = target
	depth := 0
	pc := 0

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic at pc %v - %v", pc, r)
		}
	}()

	for pc = 0; pc < len(program); pc++ {
		inst := program[pc]
		frame := &stack[depth]
		fmt.Printf("PC: %v Op: %v frame: %v\n", pc, inst, frame)
		switch inst.Op {
		case Read:
			switch inst.Type {
			case Null:
				break
			case Boolean:
				frame.Boolean, err = readBool(r)
				break
			case Int:
				frame.Int, err = readInt(r)
				break
			case Long:
				frame.Long, err = readLong(r)
				break
			case Float:
				frame.Float, err = readFloat(r)
				break
			case Double:
				frame.Double, err = readDouble(r)
				break
			case Bytes:
				frame.Bytes, err = readBytes(r)
				break
			case String:
				frame.String, err = readString(r)
				break
			case MapKey:
				frame.MapKey, err = readString(r)
				break
			case Length:
				frame.Length, err = readLong(r)
				break
			case UnionElem:
				frame.UnionElem, err = readLong(r)
				break
			}
			break
		case Set:
			switch inst.Type {
			case Null:
				break
			case Boolean:
				frame.Target.SetBoolean(inst.Field, frame.Boolean)
				break
			case Int:
				frame.Target.SetInt(inst.Field, frame.Int)
				break
			case Long:
				frame.Target.SetLong(inst.Field, frame.Long)
				break
			case Float:
				frame.Target.SetFloat(inst.Field, frame.Float)
				break
			case Double:
				frame.Target.SetDouble(inst.Field, frame.Double)
				break
			case Bytes:
				frame.Target.SetBytes(inst.Field, frame.Bytes)
				break
			case String:
				frame.Target.SetString(inst.Field, frame.String)
				break
			}
			break
		case Init:
			frame.Target.Init(inst.Field)
			break
		case Enter:
			depth += 1
			stack[depth].Target = frame.Target.Get(inst.Field)
			break
		case Exit:
			depth -= 1
			break
		case BlockStart:
			// If we're starting a block, read the header
			if frame.Length == 0 {
				stack[depth].BlockStart = pc
				frame.Length, err = readLong(r)
				// If the header is 0, the array/map is over
				if frame.Length == 0 {
					for program[pc].Op != BlockEnd {
						pc += 1
					}
				}
			}
			frame.Length -= 1
			break
		case BlockEnd:
			pc = stack[depth].BlockStart - 1
			break
		}

		if err != nil {
			return err
		}
	}
	return nil
}
