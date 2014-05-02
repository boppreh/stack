package core

import (
	"regexp"
)

func sMatches(i In, o Out) {
	pattern := "^" + i().(string) + "$"
	text := i().(string)
	result, _ := regexp.MatchString(pattern, text)
	o(result)
}
