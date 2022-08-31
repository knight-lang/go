package knight

import (
	"fmt"
)

// Ast represents a Knight code.
type Ast struct {
	fun  *Function
	args []Value
}

// Compile-time assertion that `Ast`s implements the `Value` interface.
var _ Value = &Ast{}

// NewAst creates an `Ast`.
func NewAst(fun *Function, args []Value) *Ast {
	if len(args) != fun.arity {
		panic(fmt.Sprintf("arg len mismatch (expected %d for %q, got %d)",
			fun.arity, fun.name, len(args)))
	}
	return &Ast{fun: fun, args: args}
}

// Run passes `a`'s arguments to `a`'s function.
func (a *Ast) Run() (Value, error) {
	return a.fun.fn(a.args)
}

// Dump writes a debugging representation of `a` to stdout.
func (a *Ast) Dump() {
	fmt.Printf("Ast(%c", a.fun.name)

	for _, arg := range a.args {
		fmt.Print(", ")
		arg.Dump()
	}

	fmt.Print(")")
}
