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
	var pc int
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic at pc %v - %v", pc, r)
		}
	}()

	return evalInner(r, program, target, &pc)
}

func evalInner(r io.Reader, program *Program, target types.Field, pc *int) (err error) {
	var loop int64

	frame := stackFrame{}
	for ; *pc < len(program.Instructions); *pc++ {
		inst := program.Instructions[*pc]
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
				target.SetBoolean(frame.Boolean)
				break
			case Int:
				target.SetInt(frame.Int)
				break
			case Long:
				target.SetLong(frame.Long)
				break
			case Float:
				target.SetFloat(frame.Float)
				break
			case Double:
				target.SetDouble(frame.Double)
				break
			case Bytes:
				target.SetBytes(frame.Bytes)
				break
			case String:
				target.SetString(frame.String)
				break
			}
			break
		case SetDefault:
			target.SetDefault(inst.Operand)
			break
		case Enter:
			*pc += 1
			if err = evalInner(r, program, target.Get(inst.Operand), pc); err != nil {
				return err
			}
			break
		case Exit:
			target.Finalize()
			return nil
		case AppendArray:
			*pc += 1
			if err = evalInner(r, program, target.AppendArray(), pc); err != nil {
				return err
			}
			break
		case AppendMap:
			*pc += 1
			if err = evalInner(r, program, target.AppendMap(frame.String), pc); err != nil {
				return err
			}
			break
		case Call:
			curr := *pc
			*pc = inst.Operand
			if err = evalInner(r, program, target, pc); err != nil {
				return err
			}
			*pc = curr
			break
		case Return:
			return nil
		case Jump:
			*pc = inst.Operand - 1
			break
		case EvalGreater:
			frame.Condition = (frame.Long > int64(inst.Operand))
			break
		case EvalEqual:
			frame.Condition = (frame.Long == int64(inst.Operand))
			break
		case CondJump:
			if frame.Condition {
				*pc = inst.Operand - 1
			}
			break
		case AddLong:
			frame.Long += int64(inst.Operand)
			break
		case MultLong:
			frame.Long *= int64(inst.Operand)
			break
		case PushLoop:
			loop = frame.Long
			*pc += 1
			if err = evalInner(r, program, target, pc); err != nil {
				return err
			}
			frame.Long = loop
			break
		case PopLoop:
			return nil
		case Halt:
			if inst.Operand == 0 {
				return nil
			}
			return fmt.Errorf("Runtime error: %v, frame: %v, pc: %v", program.Errors[inst.Operand-1], frame, pc)
		default:
			return fmt.Errorf("Unknown instruction %v", program.Instructions[*pc])
		}

		if err != nil {
			return err
		}
	}
	return nil
}
