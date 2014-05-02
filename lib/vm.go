package lib

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type In func() Value
type Out func(Value)

func sAdd(i In, o Out) { o(i().(int) + i().(int)) }
func sSub(i In, o Out) { o(i().(int) - i().(int)) }
func sDiv(i In, o Out) { o(i().(int) / i().(int)) }
func sMul(i In, o Out) { o(i().(int) * i().(int)) }

func sMap(i In, o Out) {
	fnList := i().([]Value)
	list := i().([]Value)

	fullFnList := make([]Value, len(fnList)+1)
	copy(fullFnList[1:], fnList)

	for i := range list {
		fullFnList[0] = list[i]
		result, _ := Run(fullFnList)
		list[i] = result[0]
	}

	o(list)
}

func sRun(i In, o Out) {
	code := i().([]Value)
	result, _ := Run(code)
	for _, value := range result {
		o(value)
	}
}

func sIf(i In, o Out) {
	value := i()
	else_ := i().([]Value)
	then := i().([]Value)

	var condition bool
	switch value.(type) {
	case int:
		condition = value != 0
	case string:
		condition = value != ""
	default:
		fmt.Println("Ops:", value)
	}

	if condition {
		o(then)
	} else {
		o(else_)
	}

	sRun(i, o)
}

func sPrint(i In, o Out) {
	v := i()
	fmt.Println(v)
	o(v)
}

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
		case func(In, Out):
			value.(func(In, Out))(s.Pop, s.Push)

		default:
			s.Push(value)
		}
	}

	return s.Exhaust(), nil
}
