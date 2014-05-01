package stack

import "errors"

type Value int

type node struct {
	value Value	
	next *node
}

type Stack struct {
	top *node
	size int
}

type Param func() Value

type Op func(Param) (Value, error)

func (s *Stack) Push(vs ...Value) {
	for _, value := range vs {
		s.top = &node{value, s.top}
		s.size++
	}
}

func (s *Stack) Pop() (v Value) {
	if s.size == 0 {
		panic("Tried to pop from an empty stack.")
	}

	v, s.top = s.top.value, s.top.next
	s.size--
	return
}

func (s *Stack) Empty() bool {
	return s.top == nil
}

func (s *Stack) Apply(op Op) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Not enough values in the stack to apply operator.")
		}
	}()

	result, op_err := op(s.Pop)
	if op_err != nil {
		return op_err
	}
	s.Push(result)

	return
}
