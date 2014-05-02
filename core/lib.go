package core

import (
	"regexp"
	"os"
	"bufio"
	"io/ioutil"
	"net/http"
	"strings"
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

func sReplace(i In, o Out) {
	replacement := i().(string)
	pattern := i().(string)
	text := i().(string)
	regex := regexp.MustCompile(pattern)
	o(regex.ReplaceAllString(text, replacement))
}

func sInput(i In, o Out) {
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		panic(err)
	}
	o(string(line))
}

func sRead(i In, o Out) {
	path := i().(string)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	o(string(bytes))
}

func sWrite(i In, o Out) {
	path := i().(string)
	contents := i().(string)
	err := ioutil.WriteFile(path, []byte(contents), 0777)
	if err != nil {
		panic(err)
	}
}

func sDelete(i In, o Out) {
	path := i().(string)
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func sGet(i In, o Out) {
	url := i().(string)
	if !strings.Contains(url, "://") {
		url = "http://" + url
	}

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	o(string(content))
}
