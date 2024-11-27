package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/knight-lang/go/knight"
)

// printAndExit 
func printAndExit(isErr bool, fmtstr string, rest ...interface{}) {
	var out *os.File

	if isErr {
		out = os.Stderr
	} else {
		out = os.Stdout
	}


	fmt.Fprintf(out, fmtstr, rest...)
	fmt.Fprint(out, "\n")

	if isErr {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Prints the usage and exits. isErr should be 
func usage(isErr bool) {
	printAndExit(isErr, "usage: %s [-h] (-e 'expr' | -f file)", os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		usage(true)
	}

	var program string
	switch os.Args[1] {
	case "-e":
		program = os.Args[2]

	case "-f":
		programBytes, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			printAndExit(true, "[FATAL] Couldn't read file contents: %s", err)
		}
		program = string(programBytes)

	case "-h":
		usage(false)

	default:
		usage(true)
	}

	if _, err := knight.Play(program); err != nil {
		printAndExit(true, "%s", err)
	}
}
