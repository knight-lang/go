package knight

import (
	"errors"
	"fmt"
)

// FnCall represents a function call (eg `+ 1 2`) in Knight. It implements Value, but
// unconditionally raise errors for all the conversion methods (as they're undefined in the specs
// for `Value`s).
type FnCall struct {
	function  *Function
	arguments []Value
}

// Compile-time assertion that FnCall implements the Value interface.
var _ Value = &FnCall{}

// NewFnCall constructs a new FnCall. It'll panic if the amount of arguments given isn't
// equal to the arity of the function.
func NewFnCall(function *Function, arguments []Value) *FnCall {
	if function.arity != len(arguments) {
		panic(fmt.Sprint("<INTERNAL BUG> function arity mismatch: expected",
			function.arity, "got", len(arguments)))
	}

	return &FnCall{function: function, arguments: arguments}
}

// Execute executes the ast by passing its arguments to its function.
func (a *FnCall) Execute() (Value, error) {
	return (a.function.fn)(a.arguments)
}

// Dump writes a debugging representation of the ast to stdout.
func (a *FnCall) Dump() {
	fmt.Printf("FnCall(%c", a.function.name)

	for _, arg := range a.arguments {
		fmt.Print(", ")
		arg.Dump()
	}

	fmt.Print(")")
}

// Conversions: They always return errors, as asts cannot be converted to other types.
func (_ *FnCall) ToString() (string, error) {
	return "", errors.New("FnCall doesn't define string conversions")
}

func (_ *FnCall) ToInt() (int, error) {
	return 0, errors.New("FnCall doesn't define int conversions")
}

func (_ *FnCall) ToBool() (bool, error) {
	return false, errors.New("FnCall doesn't define boolean conversions")
}

func (_ *FnCall) ToSlice() ([]Value, error) {
	return nil, errors.New("FnCall doesn't define list conversions")
}
