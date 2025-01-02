package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/knight-lang/go/knight"
)

// printAndExit prints out a format string and then exits with a nonzero exit status.
func printAndExit(fmtStr string, rest ...any) {
	fmt.Fprintf(os.Stderr, fmtStr, rest...)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

// usage prints the usage and exits.
func usage() {
	printAndExit("usage: %s (-e 'expr' | -f file)", os.Args[0])
}

func main() {
	// We expect exactly three arguments: The program name, `-e`/`-f`, and the expression/filename.
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
			printAndExit("[FATAL] Couldn't read file contents: %s", err)
		}
		program = string(programBytes)

	default:
		usage()
	}

	// Run the program; if there's a problem, print out the error and abort.
	if _, err := knight.Evaluate(program); err != nil {
		printAndExit("%s", err)
	}
}
