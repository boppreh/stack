package core

import (
	"errors"
	"strconv"
	"unicode"
)

type In func() Value
type Out func(Value)

func readUntil(input chan rune, delimiter rune) (text string) {
	buffer := make([]rune, 0)

	for {
		char, ok := <- input

		if char == delimiter {
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

func readWordAndCloseList(input chan rune) (text string, closedList bool) {

	buffer := make([]rune, 0)

	for {
		char, ok := <- input

		if unicode.IsSpace(char) {
			break
		}

		if char == ']' {
			closedList = true
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

	return string(buffer), closedList

}

func strToChan(sourceCode string, c chan rune) {
	for _, char := range sourceCode {
		c <- char
	}
	close(c)
}

func Parse(sourceCode string) ([]Value, error) {
	input := make(chan rune)
	go strToChan(sourceCode, input)

	s := new(Stack)
	s.Push(make([]Value, 0))

	addToken := func(token Value) {
		list := s.Pop().([]Value)
		s.Push(append(list, token))
	}

	closeList := func () {
		list := s.Pop().([]Value)
		addToken(list)
	}

	var err error

	for err == nil {
		char, ok := <-input
		if !ok {
			break
		}

		if unicode.IsNumber(char) {
			rest, closedList := readWordAndCloseList(input)

			var number int
			number, err = strconv.Atoi(string(char) + rest)
			addToken(number)

			if closedList {
				closeList()
			}
			
		} else if unicode.IsLetter(char) {
			suffix, closedList := readWordAndCloseList(input)
			name := string(char) + suffix

			function, ok := ops[name]
			if ok {
				addToken(function)
			} else {
				err = errors.New("Parser error. Unknown name " + name)
			}

			if closedList {
				closeList()
			}

		} else {
			switch char {
			case '"':
				addToken(readUntil(input, '"'))
			case '\'':
				addToken(readUntil(input, '\''))
			case ':':
				str, closedList := readWordAndCloseList(input)
				addToken(str)
				if closedList {
					closeList()
				}

			case '[':
				s.Push(make([]Value, 0))
			case ']':
				closeList()

			case '+':
				addToken(sAdd)
			case '-':
				addToken(sSub)
			case '/':
				addToken(sDiv)
			case '*':
				addToken(sMul)

			case '=':
				addToken(sEq)
			case '>':
				addToken(sGt)
			case '<':
				addToken(sLt)

			case '%':
				addToken(sMap)
			case '?':
				addToken(sIf)
			case '!':
				addToken(sIndex)
			case '.':
				addToken(sDup)

			case '$':
				addToken(sDecl)
			case '@':
				addToken(sCall)
			case '&':
				addToken(sTranspose)

			case '#':
				_ = readUntil(input, '\n')

			case ' ', '\t', '\n':
				continue

			default:
				err = errors.New("Parser error. Unexpected value " + string(char))
			}
		}
	}

	if s.size > 1 {
		err = errors.New("Parser error: literal list not closed.")
	}
	return s.Pop().([]Value), err
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
