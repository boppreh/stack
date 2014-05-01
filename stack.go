package stack

import "fmt"

type Value int

type node struct {
	value Value	
	next *node
}

type Stack struct {
	top *node
}

func (s *Stack) push(vs ...Value) {
	for _, value := range vs {
		s.top = &node{value, s.top}
	}
}

func (s *Stack) pop() (v Value) {
	v, s.top = s.top.value, s.top.next
	return
}
