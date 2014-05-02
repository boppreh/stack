package stack

type Value interface{}

type node struct {
	value Value
	next  *node
}

type Stack struct {
	top  *node
	size int
}

func New(values []Value) *Stack {
	s := new(Stack)

	for i := len(values) - 1; i >= 0; i-- {
		s.Push(values[i])
	}

	return s
}

func (s *Stack) Push(value Value) {
	s.top = &node{value, s.top}
	s.size++
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

func (s *Stack) Exhaust() (result []Value) {
	for !s.Empty() {
		result = append(result, s.Pop())
	}
	return
}
