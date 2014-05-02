package core

import (
	"testing"
)

func assertParse(t *testing.T, sourceCode string, expected... Value) {
	if expected == nil {
		expected = []Value{}
	}

	result, err := Parse(sourceCode)
	if err != nil {
		t.Errorf("%v %q", err, sourceCode)
	}
	assert(t, result, expected)
}

func assertRun(t *testing.T, sourceCode string, expected... Value) {
	if expected == nil {
		expected = []Value{}
	}

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
	assertParse(t, "")
	assertParse(t, "  ")
	assertParse(t, "  \t \n  ")

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
	assertParse(t, ":asd\\]f", "asd]f")

	assertParse(t, "[]", []Value{})
	assertParse(t, "[1]", []Value{1})
	assertParse(t, "[:name]", []Value{"name"})
	assertParse(t, "[1 2 3]", []Value{1, 2, 3})
	assertParse(t, "[1 2 3 :asdf]", []Value{1, 2, 3, "asdf"})
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

	assertRun(t, "[1 2 3] len", 3)
	assertRun(t, "'asd' len", 3)

	assertRun(t, "[1 2 3] 1 index", 2)
	assertRun(t, "[1 2 3] 1 !", 2)
	assertRun(t, "'asd' 1 !", "s")

	assertRun(t, "1 2 &", 1, 2)

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
	assertRun(t, "[:a $ 2 2 :a @ 3 4 :a @ :a @] :m $  [+] :m @", 11)

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

	assertRun(t, ":acccbd :(c+bb) find", nil)
	assertRun(t, ":acccbd :c+b? find", []Value{})
	assertRun(t, ":acccbd :(c+b?) find", []Value{"cccb"})

	assertRun(t, "'abc ac ab' :(aee?) findall", []Value{})
	assertRun(t, "'abc ac ab' :(abc?) findall", []Value{[]Value{"abc"}, []Value{"ab"}})

	assertRun(t, "'egg ham cheese' :egg :spam replace", "spam ham cheese")
	assertRun(t, "'egg hams cheese' :ham(s?) :spam$1 replace", "egg spams cheese")

	assertRun(t, "'file contents' :file.txt write")
	assertRun(t, ":file.txt read", "file contents")
	assertRun(t, ":file.txt delete")
}
