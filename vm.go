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
		char := <-input
		if char == delimiter {
			break
		} else if char == '\\' {
			char = <-input
		}

		text = append(text, char)
	}

	return string(text)
}

func parseValue(input chan rune, output chan Value) {
	for {
		char, ok := <- input
		if !ok {
			close(output)
			return
		}

		switch char {
		case '"':
			output <- parseString(input, '"')
		case ' ':
			continue
		}
	}
}

func Parse(sourceCode string) (program []Value) {
	input := make(chan rune)
	output := make(chan Value)

	go parseValue(input, output)

	for _, char := range sourceCode {
		input <- char
	}
	close(input)

	program = make([]Value, 0)	
	for {
		v, ok := <- output
		if !ok {
			return
		}
		program = append(program, v)
	}
}
