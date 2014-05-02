package core

import (
	"testing"
)

func assertParse(t *testing.T, sourceCode string, expected... Value) {
	result, err := Parse(sourceCode)
	if err != nil {
		t.Error(err)
	}
	assert(t, result, expected)
}

func assertRun(t *testing.T, sourceCode string, expected... Value) {
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
	assertParse(t, "\"first string\"", "first string")
	assertParse(t, ":symbol", "symbol")
	assertParse(t, ":symbol :symbol2", "symbol", "symbol2")
	assertParse(t, "\"a\" 'b' :symbol", "a", "b", "symbol")

	assertParse(t, "1", 1)
	assertParse(t, "12345", 12345)
	assertParse(t, "123 321", 123, 321)

	assertParse(t, "#Comment\n123", 123)
	assertParse(t, "123\n321 # Comment\n123", 123, 321, 123)
	assertParse(t, "123\n321 # Comment\n123", 123, 321, 123)

	assertParse(t, "[]", []Value{})
	assertParse(t, "[1 2 3]", []Value{1, 2, 3})
	assertParse(t, "[1 2 3] [] [1]", []Value{1, 2, 3}, []Value{}, []Value{1})
	assertParse(t, "[[]]", []Value{[]Value{}})
	assertParse(t, "[1 [2 3] 4]", []Value{1, []Value{2, 3}, 4})
	assertParse(t, "[[2 3]]", []Value{[]Value{2, 3}})
}

func TestRun(t *testing.T) {
	assertRun(t, "1", 1)
	assertRun(t, "1 2 3", 3, 2, 1)
	assertRun(t, "1 1 +", 2)
	assertRun(t, "1 1 + 2 *", 4)

	assertRun(t, "1 1 =", true)
	assertRun(t, "1 2 <", true)
	assertRun(t, "1 2 >", false)

	assertRun(t, "[1 2 3] [1 +] %", []Value{2, 3, 4})
	assertRun(t, "[1 2 3] [1 +] % [2 *] %", []Value{4, 6, 8})
	assertRun(t, "[1 2 3] [1 + 2 *] %", []Value{4, 6, 8})

	assertRun(t, "[10] !", 10)
	assertRun(t, "[10 20] !", 10, 20)
	assertRun(t, "[1 1 +] !", 2)

	assertRun(t, "1 [1] [2] ?", 1)
	assertRun(t, "0 [1] [2] ?", 2)
	assertRun(t, "1 [1 1 +] [2 2 +] ?", 2)
	assertRun(t, "0 [1 1 +] [2 2 +] ?", 4)
	assertRun(t, "0 0 = [1] [2] ?", 1)
	assertRun(t, "0 0 = [1] [2] ?", 1)

	assertRun(t, "1 [2] :a $", 1)
	assertRun(t, "[1 1 +] :a $ :a @", 2)

	assertRun(t, "1 .", 1, 1)
	assertRun(t, "1 2 . + +", 5)

	assertRun(t, "[. 5 < [1 + :inc @] [] ?] :inc $   1 :inc @", 5)

	assertRun(t, "[1 2 3] 4 append", []Value{1, 2, 3, 4})

	assertRun(t, "[1 1 +] eval", 2)
	assertRun(t, "5 dup", 5, 5)

	assertRun(t, ":5 int", 5)
	assertRun(t, "5 str", "5")
}

func TestLib(t *testing.T) {
	assertRun(t, ":aabb :a+b+ matches", true)
	assertRun(t, ":aabb :ab+ matches", false)

	assertRun(t, ":aabb :a+b+ contains", true)
	assertRun(t, ":aabb :ab+ contains", true)
	assertRun(t, ":aabb :ab+c contains", false)
}
