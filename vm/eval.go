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

	MapKey    string
	UnionType int64
}

func Eval(r io.Reader, program *Program, target types.Field) (err error) {
	callStack := make([]int, 256)
	callStackDepth := 0

	stack := make([]Frame, 256)
	stack[0].Target = target
	depth := 0

	pc := 0

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic at pc %v - %v", pc, r)
		}
	}()

	for pc = 0; pc < len(program.Instructions); pc++ {
		inst := program.Instructions[pc]
		frame := &stack[depth]
		log("PC: %v\tD:%v\tOp: %v", pc, depth, inst)
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
		case SwitchStart:
			// Skip to the case matching the UnionType in the frame
			switchDepth := 1
			for {
				pc += 1
				if program.Instructions[pc].Op == SwitchStart {
					switchDepth += 1
				}
				if program.Instructions[pc].Op == SwitchCase && program.Instructions[pc].Field == int(stack[depth].UnionType) && switchDepth == 1 {
					break
				}
				if program.Instructions[pc].Op == SwitchEnd {
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
				if program.Instructions[pc].Op == SwitchStart {
					switchDepth += 1
				}
				if program.Instructions[pc].Op == SwitchEnd {
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
		case Call:
			callStack[callStackDepth] = pc
			callStackDepth += 1
			pc = inst.Field - 1
		case Return:
			pc = callStack[callStackDepth-1]
			callStackDepth -= 1
		case Jump:
			pc = inst.Field - 1
		case ZeroJump:
			if frame.Long == 0 {
				pc = inst.Field
			}
		case DecrLong:
			frame.Long -= 1
		case Halt:
			if inst.Field == 0 {
				return nil
			} else {
				return fmt.Errorf("Runtime error: %v", program.Errors[inst.Field])
			}
		default:
			err = fmt.Errorf("Unknown instruction %v", program.Instructions[pc])
		}

		if err != nil {
			return err
		}
	}
	return nil
}
