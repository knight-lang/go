package knight

import (
	"fmt"
)

// Ast represents a function and its arguments in knight
type Ast struct {
	function  *Function
	arguments []Value
}

// Compile-time assertion that `Ast`s implements the `Value` interface.
var _ Value = &Ast{}

// Run executes the ast by passing its arguments to its function.
func (a *Ast) Run() (Value, error) {
	return a.function.fn(a.arguments)
}

// Dump writes a debugging representation of the ast to stdout.
func (a *Ast) Dump() {
	fmt.Printf("Ast(%c", a.function.name)

	for _, arg := range a.arguments {
		fmt.Print(", ")
		arg.Dump()
	}

	fmt.Print(")")
}
