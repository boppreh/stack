package stack

import "testing"

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
