package stack

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

func increment(p Param) (Value, error)   { return p().(int) + 1, nil }
func sum(p Param) (Value, error)         { return p().(int) + p().(int), nil }
func deepthought(p Param) (Value, error) { return 42, nil }

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
