package knight

import (
	"fmt"
)

// Ast is the `Value` that represents a function call (eg `+ 1 2`) in Knight.
//
// `Ast`s are only constructed within `Parser.Parse`.
type Ast struct {
	function  *Function
	arguments []Value
}

func NewAst(function *Function, arguments []Value) *Ast {
	if function.arity != len(arguments) {
		panic(fmt.Sprint("function arity mismatch: expected", function.arity, "got", len(arguments)))
	}

	return &Ast{function: function, arguments: arguments}
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
