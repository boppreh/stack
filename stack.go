package stack

type Value int

type node struct {
	value Value	
	next *node
}

type Stack struct {
	top *node
}

type Op func(*Stack) Value
type Param func() Value

func (s *Stack) Push(vs ...Value) {
	for _, value := range vs {
		s.top = &node{value, s.top}
	}
}

func (s *Stack) Pop() (v Value) {
	v, s.top = s.top.value, s.top.next
	return
}

func (s *Stack) Empty() bool {
	return s.top == nil
}

func (s *Stack) Apply(op Op) {
	s.Push(op(s))
}

func (s *Stack) Apply1(op func(Param) Value) {
	s.Push(op(s.Pop))
}

func (s *Stack) Apply2(op func(Param, Param) Value) {
	s.Push(op(s.Pop, s.Pop))
}
