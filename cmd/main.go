package main

import (
	"github.com/knight-lang/go"
	"fmt"
)

func main() {
	fmt.Println(knight.Text("foobar"))

	knight.NewVariable("foo").Dump() ; fmt.Println("")
	knight.NewVariable("foo").Dump() ; fmt.Println("")	

	knight.Number(-123).Dump()
}
