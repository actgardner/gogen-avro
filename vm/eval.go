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
		runtimeLog("PC: %v\tD:%v\tOp: %v", pc, depth, inst)
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
			case Fixed:
				frame.Bytes, err = readFixed(r, inst.Field)
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
			stack[depth].Target.Finalize()
			depth -= 1
			break
		case AppendArray:
			depth += 1
			stack[depth].Target = frame.Target.AppendArray()
			break
		case AppendMap:
			depth += 1
			stack[depth].Target = frame.Target.AppendMap(stack[depth-1].MapKey)
			break
		case BlockStart:
			// If we're starting a block, read the header
			if frame.Length == 0 {
				stack[depth].BlockStart = pc
				frame.Length, err = readLong(r)
				if err != nil {
					break
				}
				// If the header is 0, the array/map is over
				if frame.Length == 0 {
					for program[pc].Op != BlockEnd {
						pc += 1
					}
					continue
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
			switchDepth := 1
			for {
				pc += 1
				if program[pc].Op == SwitchStart {
					switchDepth += 1
				}
				if program[pc].Op == SwitchCase && program[pc].Field == int(stack[depth].UnionType) && switchDepth == 1 {
					break
				}
				if program[pc].Op == SwitchEnd {
					switchDepth -= 1
					if switchDepth == 0 {
						err = fmt.Errorf("No matching case in switch for %v", stack[depth].UnionType)
						break
					}
				}
			}
			break
		case SwitchCase:
			// Switch cases don't need an explicit break, skip to the end of the block
			switchDepth := 1
			for {
				if program[pc].Op == SwitchStart {
					switchDepth += 1
				}
				if program[pc].Op == SwitchEnd {
					switchDepth -= 1
					if switchDepth == 0 {
						break
					}
				}
				pc += 1
			}
			break
		case SwitchEnd:
			// The end of the last case, nothing to see here
			break
		default:
			err = fmt.Errorf("Unknown instruction %v", program[pc])
		}

		if err != nil {
			return err
		}
	}
	return nil
}
