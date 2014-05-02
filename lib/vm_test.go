package lib

import (
	"testing"
)

func assertParse(t *testing.T, sourceCode string, expected []Value) {
	result, err := Parse(sourceCode)
	if err != nil {
		t.Error(err)
	}
	assert(t, result, expected)
}

func assertRun(t *testing.T, sourceCode string, expected []Value) {
	program, err := Parse(sourceCode)
	if err != nil {
		t.Error(err)
	}
	result, err2 := Run(program)
	if err2 != nil {
		t.Error(err2)
	}
	assert(t, result, expected)
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

	assertParse(t, "#Comment", []Value{})
	assertParse(t, "\n  #Comment", []Value{})
	assertParse(t, "#Comment\n123", []Value{123})
	assertParse(t, "123\n321 # Comment\n123", []Value{123, 321, 123})
	assertParse(t, "123\n321 # Comment\n123", []Value{123, 321, 123})

	assertParse(t, "[]", []Value{[]Value{}})
	assertParse(t, "[1 2 3]", []Value{[]Value{1, 2, 3}})
	assertParse(t, "[1 2 3] [] [1]", []Value{[]Value{1, 2, 3}, []Value{}, []Value{1}})
}

func TestRun(t *testing.T) {
	assertRun(t, "1", []Value{1})
	assertRun(t, "1 2 3", []Value{3, 2, 1})
	assertRun(t, "1 1 +", []Value{2})
	assertRun(t, "1 1 + 2 *", []Value{4})

	assertRun(t, "[1 2 3] [1 +] %", []Value{[]Value{2, 3, 4}})
	assertRun(t, "[1 2 3] [1 +] % [2 *] %", []Value{[]Value{4, 6, 8}})
	assertRun(t, "[1 2 3] [1 + 2 *] %", []Value{[]Value{4, 6, 8}})

	assertRun(t, "[10] !", []Value{10})
	assertRun(t, "[10 20] !", []Value{10, 20})
	assertRun(t, "[1 1 +] !", []Value{2})

	assertRun(t, "1 [1] [2] ?", []Value{1})
	assertRun(t, "0 [1] [2] ?", []Value{2})
	assertRun(t, "1 [1 1 +] [2 2 +] ?", []Value{2})
	assertRun(t, "0 [1 1 +] [2 2 +] ?", []Value{4})

	assertRun(t, "1 [2] :a $", []Value{1})
	assertRun(t, "[1 1 +] :a $ :a @", []Value{2})
}
