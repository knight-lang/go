package knight

import (
	"errors"
	"fmt"
)

// Ast represents a function call (eg `+ 1 2`) in Knight. It implements Value, but unconditionally
// raise errors for all the conversion methods (as they're undefined in the specs for `Value`s).
type Ast struct {
	function  *Function
	arguments []Value
}

// Compile-time assertion that Ast implements the Value interface.
var _ Value = &Ast{}

// NewAst constructs a new ast. It'll panic if the amount of arguments given isn't equal to the
// arity of the function.
func NewAst(function *Function, arguments []Value) *Ast {
	if function.arity != len(arguments) {
		panic(fmt.Sprint("<INTERNAL BUG> function arity mismatch: expected",
			function.arity, "got", len(arguments)))
	}

	return &Ast{function: function, arguments: arguments}
}

// Execute executes the ast by passing its arguments to its function.
func (a *Ast) Execute() (Value, error) {
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

// Conversions: They always return errors, as asts cannot be converted to other types.
func (_ *Ast) ToString() (string, error) {
	return "", errors.New("Ast doesn't define string conversions")
}

func (_ *Ast) ToInt() (int, error) {
	return 0, errors.New("Ast doesn't define int conversions")
}

func (_ *Ast) ToBool() (bool, error) {
	return false, errors.New("Ast doesn't define boolean conversions")
}

func (_ *Ast) ToSlice() ([]Value, error) {
	return nil, errors.New("Ast doesn't define list conversions")
}
