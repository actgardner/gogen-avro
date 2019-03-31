package vm

import (
	"github.com/actgardner/gogen-avro/vm/types"
)

type intStack struct {
	stack []int
	pos   int
}

func (i *intStack) push(value int) {
	i.pos += 1
	if i.pos >= len(i.stack) {
		i.stack = append(i.stack, make([]int, len(i.stack))...)
	}
	i.stack[i.pos] = value
}

func (i *intStack) peek() int {
	return i.stack[i.pos]
}

func (i *intStack) pop() int {
	i.pos -= 1
	return i.stack[i.pos+1]
}

type fieldStack struct {
	stack []types.Field
	pos   int
}

func (i *fieldStack) push(value types.Field) {
	i.pos += 1
	if i.pos >= len(i.stack) {
		i.stack = append(i.stack, make([]types.Field, len(i.stack))...)
	}
	i.stack[i.pos] = value
}

func (i *fieldStack) peek() types.Field {
	return i.stack[i.pos]
}

func (i *fieldStack) pop() types.Field {
	i.pos -= 1
	return i.stack[i.pos+1]
}
