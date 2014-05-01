package stack

import "testing"

func assert(t *testing.T, result Value, expected Value) {
	if result != expected {
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

func assertError(t *testing.T, s *Stack, op Op) {
	if s.Apply(op) == nil {
		t.Errorf("Expected error, but op completely normally.")
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
}

func TestStackApply(t *testing.T) {
	s := new(Stack)
	increment := func (p Param) Value { return p() + 1 }

	s.Push(1)
	s.Apply(increment)
	assertStack(t, s, 2)

	s.Push(3)
	s.Push(5)
	s.Apply(func (p Param) Value { return p() + p() })
	assertStack(t, s, 8)

	s.Apply(func (p Param) Value { return 22 })
	assertStack(t, s, 22)

	assertError(t, s, increment)

	s.Push(1)
	s.Apply(increment, increment, increment)
	assertStack(t, s, 4)
}
