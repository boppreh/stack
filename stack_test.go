package stack

import "testing"

func assert(t *testing.T, result Value, expected Value) {
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func assertStack(t *testing.T, stack *Stack, expectedValues ...Value) {
	for key, value := range expectedValues {
		popped := stack.pop()
		if popped != value {
			t.Errorf("Expected value %v to be %v, got %v instead.", key, value, popped)
		}
	}
}

func TestStackStruct(t *testing.T) {
	s := new(Stack)

	s.push(10)
	assert(t, s.pop(), 10)

	s.push(20)
	s.push(15)
	s.push(1)
	assertStack(t, s, 1, 15, 20)
}
