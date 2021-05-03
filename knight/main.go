package main

import (
	"github.com/knight-lang/go"
	"fmt"
	"strings"
	"log"
	"os"
	"io/ioutil"
)

func run(s string) knight.Value {
	val, err := knight.Parse(strings.NewReader(s))

	if err != nil {
		log.Fatal(err)
	}

	val.Dump()
	fmt.Println()
	return val
}

func main() {
	if len(os.Args) != 3 || (os.Args[1] != "-e" && os.Args[1] != "-f") {
		fmt.Printf("usage: %s (-e 'expr' | -f file)", os.Args[0])
		os.Exit(1)
	}

	var program string

	if os.Args[1] == "-e" {
		program = os.Args[2]
	} else {
		program_bytes, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			log.Fatal(err)
		}

		program = string(program_bytes)
	}

	_, err := knight.Run(program)

	if err != nil {
		log.Fatal(err)
	}
}
