package stack

import "testing"

func assert(t *testing.T, result Value, expected Value) {
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func assertStack(t *testing.T, stack *Stack, expectedValues ...Value) {
	for key, value := range expectedValues {
		popped := stack.Pop()
		if popped != value {
			t.Errorf("Expected value %v to be %v, got %v instead.", key, value, popped)
		}
	}

	if !stack.Empty() {
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
}

func TestStackApply(t *testing.T) {
	s := new(Stack)
	s.Push(1)
	s.Push(1)

	s.Apply(func (s *Stack) Value {
		return s.Pop() + s.Pop()
	})

	assertStack(t, s, 2)
}

func TestStackApplyN(t *testing.T) {
	s := new(Stack)
	s.Push(1)

	s.Apply1(func (p Param) Value {
		return p() + 1
	})

	assertStack(t, s, 2)

	s.Push(1)
	s.Push(1)

	s.Apply2(func (a Param, b Param) Value {
		return a() + b()
	})

	assertStack(t, s, 2)
}
