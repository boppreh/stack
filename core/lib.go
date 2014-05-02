package core

import (
	"regexp"
)

func sMatches(i In, o Out) {
	o("^" + i().(string) + "$")
	sContains(i, o)
}

func sContains(i In, o Out) {
	pattern := i().(string)
	text := i().(string)
	result, _ := regexp.MatchString(pattern, text)
	o(result)
}
