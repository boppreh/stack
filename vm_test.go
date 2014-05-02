package stack

import (
	"reflect"
	"testing"
)

func assertParse(t *testing.T, sourceCode string, expected []Value) {
	result := Parse(sourceCode)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v.", expected, result)
	}
}

func TestParse(t *testing.T) {
	assertParse(t, "", []Value{})
	assertParse(t, "   ", []Value{})
	assertParse(t, " \n  ", []Value{})

	assertParse(t, "\"first string\"", []Value{"first string"})
	assertParse(t, ":symbol", []Value{"symbol"})
	assertParse(t, ":symbol :symbol2", []Value{"symbol", "symbol2"})
	assertParse(t, "\"a\" 'b' :symbol", []Value{"a", "b", "symbol"})

	assertParse(t, "1", []Value{1})
	assertParse(t, "12345", []Value{12345})
	assertParse(t, "123 321", []Value{123, 321})

	assertParse(t, "123\n321 # Comment\n123", []Value{123, 321, 123})
}
