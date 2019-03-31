// The GADGT VM implementation and instruction set
package vm

import (
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/vm/types"
)

type stackFrame struct {
	Boolean   bool
	Int       int32
	Long      int64
	Float     float32
	Double    float64
	Bytes     []byte
	String    string
	Condition bool
}

func Eval(r io.Reader, program *Program, target types.Field) (err error) {
	// Stack of pointers for returning from function calls
	callStack := &intStack{stack: make([]int, 8), pos: -1}

	// Stack of loop variables
	loopStack := &intStack{stack: make([]int, 8), pos: -1}

	// Stack of target Fields for assigning values
	targetStack := &fieldStack{stack: make([]types.Field, 8), pos: -1}
	targetStack.push(target)

	frame := stackFrame{}

	pc := 0

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic at pc %v - %v", pc, r)
		}
	}()

	for pc = 0; pc < len(program.Instructions); pc++ {
		inst := program.Instructions[pc]
		switch inst.Op {
		case Read:
			switch inst.Operand {
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
			case UnusedLong:
				_, err = readLong(r)
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
			default:
				frame.Bytes, err = readFixed(r, inst.Operand-11)
				break
			}
			break
		case Set:
			switch inst.Operand {
			case Null:
				break
			case Boolean:
				targetStack.peek().SetBoolean(frame.Boolean)
				break
			case Int:
				targetStack.peek().SetInt(frame.Int)
				break
			case Long:
				targetStack.peek().SetLong(frame.Long)
				break
			case Float:
				targetStack.peek().SetFloat(frame.Float)
				break
			case Double:
				targetStack.peek().SetDouble(frame.Double)
				break
			case Bytes:
				targetStack.peek().SetBytes(frame.Bytes)
				break
			case String:
				targetStack.peek().SetString(frame.String)
				break
			}
			break
		case SetDefault:
			targetStack.peek().SetDefault(inst.Operand)
			break
		case Enter:
			targetStack.push(targetStack.peek().Get(inst.Operand))
			break
		case Exit:
			targetStack.pop().Finalize()
			break
		case AppendArray:
			targetStack.push(targetStack.peek().AppendArray())
			break
		case AppendMap:
			targetStack.push(targetStack.peek().AppendMap(frame.String))
			break
		case Call:
			callStack.push(pc)
			pc = inst.Operand - 1
			break
		case Return:
			pc = callStack.pop()
			break
		case Jump:
			pc = inst.Operand - 1
			break
		case EvalGreater:
			frame.Condition = (frame.Long > int64(inst.Operand))
			break
		case EvalEqual:
			frame.Condition = (frame.Long == int64(inst.Operand))
			break
		case CondJump:
			if frame.Condition {
				pc = inst.Operand - 1
			}
			break
		case AddLong:
			frame.Long += int64(inst.Operand)
			break
		case MultLong:
			frame.Long *= int64(inst.Operand)
			break
		case PushLoop:
			loopStack.push(int(frame.Long))
			break
		case PopLoop:
			frame.Long = int64(loopStack.pop())
			break
		case Halt:
			if inst.Operand == 0 {
				return nil
			}
			return fmt.Errorf("Runtime error: %v, frame: %v, pc: %v", program.Errors[inst.Operand-1], frame, pc)
		default:
			return fmt.Errorf("Unknown instruction %v", program.Instructions[pc])
		}

		if err != nil {
			return err
		}
	}
	return nil
}
