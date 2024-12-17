package knight

import (
	"fmt"
)

// Ast represents a function call (eg `+ 1 2`) in Knight.
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
		panic(fmt.Sprint("[BUG] function arity mismatch: expected", function.arity, "got", len(arguments)))
	}

	return &Ast{function: function, arguments: arguments}
}

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

// ToString converts the result of running the Ast to a String.
func (a *Ast) ToString() (String, error) {
	ran, err := a.Run()
	if err != nil {
		return "", err
	}

	return ran.ToString()
}

// ToInteger converts the result of running the Ast to an Integer.
func (a *Ast) ToInteger() (Integer, error) {
	ran, err := a.Run()
	if err != nil {
		return 0, err
	}

	return ran.ToInteger()
}

// ToBoolean converts the result of running the Ast to a Boolean.
func (a *Ast) ToBoolean() (Boolean, error) {
	ran, err := a.Run()
	if err != nil {
		return false, err
	}

	return ran.ToBoolean()
}

// ToList converts the result of running the Ast to a List.
func (a *Ast) ToList() (List, error) {
	ran, err := a.Run()
	if err != nil {
		return nil, err
	}

	return ran.ToList()
}
