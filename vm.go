package stack

import (
	"strconv"
	"unicode"
)

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


	for {
		var token Value
		char, ok := <- input
		if !ok {
			return
		}

		if unicode.IsNumber(char) {
			token = parseNumber(input, char)
		}

		switch char {
		case '"':
			token = parseString(input, '"')
		case '\'':
			token = parseString(input, '\'')
		case ':':
			token = parseString(input, ' ')

		case '#':
			ignoreComment(input)
			continue

		case ' ', '\t', '\n':
			continue
		}
		program = append(program, token)
	}

	return
}
