package lib

import (
	"fmt"
	"errors"
	"strconv"
	"unicode"
)

type In func() Value
type Out func(Value)

func ignoreComment(input chan rune) {
	for {
		char, ok := <-input
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

func strToChan(sourceCode string, c chan rune) {
	for _, char := range sourceCode {
		c <- char
	}
	close(c)
}

func parseChan(input chan rune) (program []Value, err error) {
	program = make([]Value, 0)

	for {
		var token Value
		char, ok := <-input
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

			case '+':
				token = sAdd
			case '-':
				token = sSub
			case '/':
				token = sDiv
			case '*':
				token = sMul

			case '%':
				token = sMap
			case '?':
				token = sIf
			case '!':
				token = sRun

			case '$':
				token = sDecl
			case '@':
				token = sCall

			case '#':
				ignoreComment(input)
				continue

			case ' ', '\t', '\n':
				continue

			default:
				name := string(char) + parseString(input, ' ')
				token, ok = ops[name]
				fmt.Println(name, token, ok)
				if !ok || token == nil {
					return nil, errors.New("Parser error. Unexpected value " + name)
				}
			}
		}

		program = append(program, token)
	}

	return program, nil
}

func Parse(sourceCode string) (program []Value, err error) {
	input := make(chan rune)
	go strToChan(sourceCode, input)
	program, err = parseChan(input)
	return
}

func Run(program []Value) ([]Value, error) {
	s := new(Stack)

	for _, value := range program {
		switch value.(type) {
		case func(In, Out):
			value.(func(In, Out))(s.Pop, s.Push)

		default:
			s.Push(value)
		}
	}

	return s.Exhaust(), nil
}
