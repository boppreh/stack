package core

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, result Value, expected Value) {
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func assertStack(t *testing.T, s *Stack, expectedValues ...Value) {
	for key, value := range expectedValues {
		popped := s.Pop()
		if popped != value {
			t.Errorf("Expected value %v to be %v, got %v instead.", key, value, popped)
		}
	}

	if !s.Empty() {
		t.Errorf("Stack has more elements than expected.")
	}
}

func TestStackStruct(t *testing.T) {
	s := new(Stack)

	s.Push(10)
	assert(t, s.Pop(), 10)

	s.Push(20)
	s.Push(15)
	s.Push(1)
	assertStack(t, s, 1, 15, 20)

	assertStack(t, New([]Value{1, 15, 20}), 1, 15, 20)
}
