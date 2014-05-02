package stack

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func add(p Param) (Value, error) { return p().(int) + p().(int), nil }
func sub(p Param) (Value, error) { return p().(int) - p().(int), nil }
func div(p Param) (Value, error) { return p().(int) / p().(int), nil }
func mul(p Param) (Value, error) { return p().(int) * p().(int), nil }

func ignoreComment(input chan rune) {
	for {
		char, ok := <- input
		if char == '\n' || !ok {
			return
		}
	}
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

func parseNumber(input chan rune, firstDigit rune) int {
	number := make([]rune, 0)
	number = append(number, firstDigit)
	for {
		digit, ok := <-input
		if !ok || !unicode.IsNumber(digit) {
			break
		}
		number = append(number, digit)
	}

	result, _ := strconv.Atoi(string(number))
	return result
}

func parseList(input chan rune, delimiter rune) []Value {
	listChan := make(chan rune)

	go func() {
		for {
			char, ok := <-input
			if char == delimiter || !ok {
				close(listChan)
				return
			}
			listChan <- char
		}
	}()

	list, _ := parseChan(listChan)
	return list
}

func lexer(sourceCode string, c chan rune) {
	for _, char := range sourceCode {
		c <- char
	}
	close(c)
}

func parseChan(input chan rune) (program []Value, err error) {
	program = make([]Value, 0)	

	for {
		var token Value
		char, ok := <- input
		if !ok {
			return
		}

		if unicode.IsNumber(char) {
			token = parseNumber(input, char)
		} else {
			switch char {
			case '"':
				token = parseString(input, '"')
			case '\'':
				token = parseString(input, '\'')
			case ':':
				token = parseString(input, ' ')

			case '[':
				token = parseList(input, ']')

			case '+': token = sum
			case '-': token = sub
			case '/': token = div
			case '*': token = mul

			case '#':
				ignoreComment(input)
				continue

			case ' ', '\t', '\n':
				continue

			default:
				return nil, errors.New("Parse error")
				continue
			}
		}

		program = append(program, token)
	}

	return program, nil
}

func Parse(sourceCode string) (program []Value, err error) {
	input := make(chan rune)
	go lexer(sourceCode, input)
	program, err = parseChan(input)
	return
}

func Run(program []Value) ([]Value, error) {
	s := new(Stack)

	for _, value := range program {
		switch value.(type) {
		case func(Param) (Value, error):
			s.Push(value.(func(Param) (Value, error))(s.Pop))

		default:
			s.Push(value)
		}
	}

	return s.Exhaust(), nil
}
