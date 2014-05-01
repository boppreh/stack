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

func increment(p Param) (Value, error) { return p() + 1, nil }
func sum(p Param) (Value, error) { return p() + p(), nil }
func deepthought(p Param) (Value, error) { return 42, nil }

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
	s.Apply(increment)
	assertStack(t, s, 2)

	s.Push(3)
	s.Push(5)
	s.Apply(sum)
	assertStack(t, s, 8)

	s.Apply(deepthought)
	assertStack(t, s, 42)

	assertError(t, s, increment)
}
