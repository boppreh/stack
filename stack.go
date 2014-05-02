package main

import (
	"github.com/boppreh/stack/lib"
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	var bytes []byte
	var err error

	if len(os.Args) == 1 {
		bytes, err = ioutil.ReadAll(os.Stdin)	
	} else {
		bytes, err = ioutil.ReadFile(os.Args[1])
	}

	if err != nil {
		panic(err)
	}

	program, err := lib.Parse(string(bytes))
	result, err := lib.Run(program)
	if err != nil {
		panic(err)
	}

	for _, value := range result {
		fmt.Println(value)
	}
}
