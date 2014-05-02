package core

import (
	"errors"
	"strconv"
	"unicode"
)

type In func() Value
type Out func(Value)

func readUntil(input chan rune, condition func(rune) bool) (text string) {
	buffer := make([]rune, 0)

	for {
		char, ok := <- input

		if condition(char) {
			break
		}

		if char == '\\' {
			char, ok = <- input
		}

		if !ok {
			break
		}

		buffer = append(buffer, char)
	}

	return string(buffer)
}

func parseString(input chan rune, delimiter rune) (string, error) {
	return readUntil(input, func(char rune) bool {
		return char == delimiter
	}), nil
}

func parseNumber(input chan rune, firstDigit rune) (int, error) {
	strNumber := readUntil(input, func (char rune) bool {
		return !unicode.IsNumber(char)
	})

	return strconv.Atoi(string(firstDigit) + strNumber)
}

func parseList(input chan rune, delimiter rune) (list []Value, err error) {
	list, err = parseChan(input)

	if err == nil {
		err = errors.New("Parser error: list not closed.")
	} else if err.Error() == "Unexpected ]." {
		err = nil
	}

	return
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
		char, ok := <-input
		if !ok {
			return
		}

		var token Value
		err = nil

		if unicode.IsNumber(char) {
			token, err = parseNumber(input, char)
			
		} else if unicode.IsLetter(char) {
			name := string(char)

			var nameSuffix string
			nameSuffix, err = parseString(input, ' ')
			if err != nil {
				return
			}

			name += nameSuffix

			if name[len(name) - 1] == ']' {
				name = name[:len(name) - 1]
				err = errors.New("Unexpected ].")
			}

			token, ok = ops[name]
			if !ok {
				err = errors.New("Parser error. Unknown name " + name)
				return
			}

		} else {
			switch char {
			case '"':
				token, err = parseString(input, '"')
			case '\'':
				token, err = parseString(input, '\'')
			case ':':
				token, err = parseString(input, ' ')

			case '[':
				token, err = parseList(input, ']')
			case ']':
				err = errors.New("Unexpected ].")

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
				_, err = parseString(input, '\n')
				token = nil

			case ' ', '\t', '\n':
				token = nil

			default:
				err = errors.New("Parser error. Unexpected value " + string(char))
				return
			}
		}
		
		if token != nil {
			program = append(program, token)
		}

		if err != nil {
			return
		}
	}

	return
}

func Parse(sourceCode string) ([]Value, error) {
	input := make(chan rune)
	go strToChan(sourceCode, input)
	return parseChan(input)
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
