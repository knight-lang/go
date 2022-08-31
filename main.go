package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/knight-lang/go"
)

// import "net/http"
// import _ "net/http/pprof"

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	if len(os.Args) != 3 || (os.Args[1] != "-e" && os.Args[1] != "-f") {
		fmt.Printf("usage: %s (-e 'expr' | -f file)", os.Args[0])
		os.Exit(1)
	}

	var program string

	if os.Args[1] == "-e" {
		program = os.Args[2]
	} else {
		programBytes, err := ioutil.ReadFile(os.Args[2])

		if err != nil {
			log.Fatal(err)
		}

		program = string(programBytes)
	}

	_, err := knight.Play(program)

	if err != nil {
		log.Fatal(err)
	}
}
