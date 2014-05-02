package core

import (
	"regexp"
	"fmt"
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

func sFind(i In, o Out) {
	regex := regexp.MustCompile(i().(string))
	fmt.Println(regex)
	submatches := regex.FindStringSubmatch(i().(string))
	if submatches == nil {
		o(nil)
	} else {
		o(submatches[1:])
	}
}

func sFindAll(i In, o Out) {
	regex := regexp.MustCompile(i().(string))
	allSubmatches := regex.FindAllStringSubmatch(i().(string), -1)
	fmt.Println(allSubmatches)
	for i, submatch := range allSubmatches {
		allSubmatches[i] = submatch[1:]
	}
	o(allSubmatches)
}
