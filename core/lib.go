package core

import (
	"regexp"
)

func strToValue(strings []string) []Value {
	values := make([]Value, 0)
	for _, str := range strings {
		values = append(values, str)
	}
	return values
}

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

func sFind(i In, o Out) {
	regex := regexp.MustCompile(i().(string))
	submatches := regex.FindStringSubmatch(i().(string))
	if submatches == nil {
		o(nil)
	} else {
		o(strToValue(submatches[1:]))
	}
}

func sFindAll(i In, o Out) {
	regex := regexp.MustCompile(i().(string))
	allSubmatches := regex.FindAllStringSubmatch(i().(string), -1)

	if allSubmatches == nil {
		o([]Value{})
		return
	}

	results := make([]Value, len(allSubmatches))

	for i, submatch := range allSubmatches {
		results[i] = strToValue(submatch[1:])
	}

	o(results)
}
