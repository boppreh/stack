package main

import (
	"github.com/boppreh/stack/lib"
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
)

func runAndPrint(sourceCode string) {
	program, err := lib.Parse(sourceCode)
	if err != nil {
		fmt.Println(err)
	}

	result, err := lib.Run(program)
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range result {
		fmt.Println(value)
	}
}

func main() {
	if len(os.Args) == 1 {
		for {
			fmt.Print(">> ")
			bio := bufio.NewReader(os.Stdin)
			line, _, err := bio.ReadLine()
			if err != nil {
				panic(err)
			}
			runAndPrint(string(line))
		}
	} else {
		bytes, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		runAndPrint(string(bytes))
	}
}
