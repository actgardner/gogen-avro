package vm

type Op int

const (
	Read Op = iota
	Set
	Enter
	Exit
	AppendArray
	AppendMap
	SwitchStart
	SwitchCase
	SwitchEnd
	Call
	Return
	Halt
	Jump
	ZeroJump
	DecrLong
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
	case SwitchStart:
		return "switch_start"
	case SwitchCase:
		return "switch_case"
	case SwitchEnd:
		return "switch_end"
	case Call:
		return "call"
	case Return:
		return "return"
	case Halt:
		return "halt"
	case Jump:
		return "jump"
	case ZeroJump:
		return "zero_jump"
	case DecrLong:
		return "decr_long"
	}
	return "Unknown"
}
