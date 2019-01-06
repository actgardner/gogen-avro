package vm

import (
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/types"
)

type Frame struct {
	Target types.Field

	Boolean bool
	Int     int32
	Long    int64
	Float   float32
	Double  float64
	Bytes   []byte
	String  string

	MapKey string
	Length int64

	BlockStart int
	UnionType  int64
}

func Eval(r io.Reader, program []Instruction, target types.Field) (err error) {
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
		//fmt.Printf("PC: %v Op: %v frame: %v\n", pc, inst, frame)
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
				frame.UnionType, err = readLong(r)
				break
			}
			break
		case Set:
			switch inst.Type {
			case Null:
				break
			case Boolean:
				frame.Target.SetBoolean(frame.Boolean)
				break
			case Int:
				frame.Target.SetInt(frame.Int)
				break
			case Long:
				frame.Target.SetLong(frame.Long)
				break
			case Float:
				frame.Target.SetFloat(frame.Float)
				break
			case Double:
				frame.Target.SetDouble(frame.Double)
				break
			case Bytes:
				frame.Target.SetBytes(frame.Bytes)
				break
			case String:
				frame.Target.SetString(frame.String)
				break
			case UnionElem:
				frame.Target.SetUnionElem(frame.UnionType)
			}
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
			// Loop back to the beginning of the loop
			pc = stack[depth].BlockStart - 1
			break
		case SwitchStart:
			// Skip to the case matching the UnionType in the frame
			for {
				if program[pc].Op == SwitchCase && program[pc].Field == int(stack[depth].UnionType) {
					break
				}
				if program[pc].Op == SwitchEnd {
					err = fmt.Errorf("No matching case in switch for %v", stack[depth].UnionType)
				}
				pc += 1
			}
			break
		case SwitchCase:
			// Switch cases don't need an explicit break, skip to the end of the block
			for {
				if program[pc].Op == SwitchEnd {
					break
				}
				pc += 1
			}
			break
		case SwitchEnd:
			// The end of the last case, nothing to see here
			break
		}

		if err != nil {
			return err
		}
	}
	return nil
}
