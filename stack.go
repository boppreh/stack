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

func (s *Stack) apply(op Op) {
	s.Push(op(s))
}
