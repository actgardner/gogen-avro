package vm

// OP represents an opcode for the VM. Operations take 0 or 1 operands.
type Op int

const (
	// Read a value of the operand type from the wire and put itin the frame
	Read Op = iota

	// Set the current target to the value of the operand type from the frame
	Set

	// Allocate a new frame and make the target the field with the operand index
	Enter

	// Move to the previous frame
	Exit

	// Append a value to the current target and enter the new value
	AppendArray

	// Append a new key-value pair (where the key is the String value in the current frame) to the current target and enter the new value
	AppendMap

	// Push the current address onto the call stack and move the PC to the operand address
	Call

	// Pop the top value frmm the call stack and set the PC to that address
	Return

	// Stop the VM. If the operand is greater than zero, look up the corresponding error message and return it
	Halt

	// Move the PC to the operand
	Jump

	// Evaluate whether the Long register is equal to the operand, and set the condition register to the result
	EvalEqual

	// Evaluate whether the Long register is greater than the operand, and set the condition register to the result
	EvalGreater

	// If the condition register is true, evaluate the next instruction. Otherwise skip to the following instruction
	Cond

	// Add the operand value to the Long register
	AddLong

	// Multiply the operand value by the Long register
	MultLong
)

func (o Op) String() string {
	switch o {
	case Read:
		return "read"
	case Set:
		return "set"
	case Enter:
		return "enter"
	case Exit:
		return "exit"
	case AppendArray:
		return "append_array"
	case AppendMap:
		return "append_map"
	case Call:
		return "call"
	case Return:
		return "return"
	case Halt:
		return "halt"
	case Jump:
		return "jump"
	case EvalEqual:
		return "eval_equal"
	case EvalLess:
		return "eval_less"
	case Cond:
		return "cond"
	case AddLong:
		return "add_long"
	case MultLong:
		return "mult_long"
	}
	return "Unknown"
}
