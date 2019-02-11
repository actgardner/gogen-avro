package vm

type Op int

const (
	Read Op = iota
	Set
	Enter
	Exit
	AppendArray
	AppendMap
	Call
	Return
	Halt
	Jump
	CondJump
	AddLong
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
	case CondJump:
		return "cond_jump"
	case AddLong:
		return "add_long"
	}
	return "Unknown"
}
