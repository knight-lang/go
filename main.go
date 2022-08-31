package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/knight-lang/go"
)

func die(fmtstr string, rest ...interface{}) {
	fmt.Fprintf(os.Stderr, fmtstr, rest...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func usage() {
	die("usage: %s (-e 'expr' | -f file)", os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	var program string
	switch os.Args[1] {
	case "-e":
		program = os.Args[2]

	case "-f":
		programBytes, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			die("couldn't read file contents: %s", err)
		}
		program = string(programBytes)

	default:
		usage()
	}

	if _, err := knight.Play(program); err != nil {
		die("%s", err)
	}
}
