package core

import (
	"errors"
	"strconv"
	"unicode"
)

type In func() Value
type Out func(Value)

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
	list, err := parseChan(input)
	if err == nil || err.Error() != "Unexpected ]." {
		panic("Parser error: list not closed.")
	}
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

	var closedBracket bool
	for {
		if closedBracket {
			return program, errors.New("Unexpected ].")
		}

		var token Value
		char, ok := <-input
		if !ok {
			return
		}

		if unicode.IsNumber(char) {
			token = parseNumber(input, char)
		} else if unicode.IsLetter(char) {
			name := string(char) + parseString(input, ' ')

			if name[len(name) - 1] == ']' {
				name = name[:len(name) - 1]
				closedBracket = true
			}

			token, ok = ops[name]
			if !ok || token == nil {
				return nil, errors.New("Parser error. Unexpected name " + name)
			}
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
			case ']':
				return program, errors.New("Unexpected ].")

			case '+':
				token = sAdd
			case '-':
				token = sSub
			case '/':
				token = sDiv
			case '*':
				token = sMul

			case '=':
				token = sEq
			case '>':
				token = sGt
			case '<':
				token = sLt

			case '%':
				token = sMap
			case '?':
				token = sIf
			case '!':
				token = sEval
			case '.':
				token = sDup

			case '$':
				token = sDecl
			case '@':
				token = sCall

			case '#':
				_ = parseString(input, '\n')
				continue

			case ' ', '\t', '\n':
				continue

			default:
				return nil, errors.New("Parser error. Unexpected value " + string(char))
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

func pushAndRun(i In, o Out, valuesToAdd []Value) {
	for _, value := range valuesToAdd {
		switch value.(type) {
		case func(In, Out):
			value.(func(In, Out))(i, o)

		default:
			o(value)
		}
	}
}

func Run(program []Value) ([]Value, error) {
	s := new(Stack)
	pushAndRun(s.Pop, s.Push, program)
	return s.Exhaust(), nil
}
