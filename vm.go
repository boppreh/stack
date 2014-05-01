package stack

import (
	"fmt"
)

func RunOps(inputs []Value, ops []Op) ([]Value, error) {
	s := New(inputs)

	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			return nil, err
		}
	}

	return s.Exhaust(), nil
}

type OpCode int
type Program []OpCode

func Run(program Program) ([]Value, error) {
	s := new(Stack)

	r := s.Push
	p := s.Pop

	for _, opCode := range program {
		switch opCode {
		case 0:
			r(0)
		case 1:
			r(p().(int) + 1)
		case 2:
			r(p().(int) + p().(int))

		case 10:
			fmt.Print(p())
		}
	}

	return s.Exhaust(), nil
}

func parseString(input chan rune, delimiter rune) string {
	text := make([]rune, 0)

	for {
		char, ok := <-input
		if char == delimiter || !ok {
			break
		} else if char == '\\' {
			char = <-input
		}

		text = append(text, char)
	}

	return string(text)
}

func lexer(sourceCode string, c chan rune) {
	for _, char := range sourceCode {
		c <- char
	}
	close(c)
}

func Parse(sourceCode string) (program []Value) {
	input := make(chan rune)
	go lexer(sourceCode, input)

	program = make([]Value, 0)	

	var token Value

	for {
		char, ok := <- input
		if !ok {
			return
		}

		switch char {
		case '"':
			token = parseString(input, '"')
		case '\'':
			token = parseString(input, '\'')
		case ':':
			token = parseString(input, ' ')
		case ' ':
			continue
		}
		program = append(program, token)
	}

	return
}
