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

func assertResult(t *testing.T, inputs []Value, ops []Op, expected []Value) {
	results, err := Run(inputs, ops)

	if err != nil {
		t.Error(err)
	}

	if len(results) != len(expected) {
		t.Errorf("Expected %v results, got %v.", len(expected), len(results))
	}

	for i := range results {
		assert(t, results[i], expected[i])
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

func TestRun(t *testing.T) {
	_, err := Run([]Value{}, []Op{increment})
	if err == nil {
		t.Errorf("Insufficient values should raise an error.")
	}

	assertResult(t, []Value{}, []Op{}, []Value{})
	assertResult(t, []Value{2, 2}, []Op{sum}, []Value{4})
	assertResult(t, []Value{2, 2, 2}, []Op{sum, sum}, []Value{6})
	assertResult(t, []Value{}, []Op{deepthought, increment}, []Value{43})
}
