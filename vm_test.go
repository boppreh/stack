package stack

import (
	"reflect"
	"testing"
)

func assertResult(t *testing.T, inputs []Value, ops []Op, expected []Value) {
	results, err := RunOps(inputs, ops)

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

func TestRunOps(t *testing.T) {
	_, err := RunOps([]Value{}, []Op{increment})
	if err == nil {
		t.Errorf("Insufficient values should raise an error.")
	}

	assertResult(t, []Value{}, []Op{}, []Value{})
	assertResult(t, []Value{2, 2}, []Op{sum}, []Value{4})
	assertResult(t, []Value{2, 2, 2}, []Op{sum, sum}, []Value{6})
	assertResult(t, []Value{}, []Op{deepthought, increment}, []Value{43})
}

func TestRun(t *testing.T) {
	result, err := Run(Program{0, 0, 0})
	if err != nil {
		t.Error(err)
	}
	expected := []Value{0, 0, 0}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v.", expected, result)
	}

	result, err = Run(Program{0, 1, 1, 0, 1, 2})
	if err != nil {
		t.Error(err)
	}
	expected = []Value{3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v.", expected, result)
	}
}

func TestParse(t *testing.T) {
	result := Parse("\"first string\"")
	expected := []Value{"first string"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v.", expected, result)
	}
	
}
