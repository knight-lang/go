package knight

import (
	"fmt"
	"math/rand"
)

type Function struct {
	name  rune
	arity int
	body  func([]Value) (Value, error)
}

type Ast struct {
	fn   *Function
	args []Value
}

func (a *Ast) Run() (Value, error) {
	return a.fn.body(a.args)
}

func (a *Ast) Dump() {
	fmt.Printf("Function(%c", a.fn.name)

	for _, arg := range a.args {
		fmt.Print(", ")
		arg.Dump()
	}

	fmt.Print(")")
}

func init() {

}

func Prompt([]Value) Value {
	return Text("A")
}

func Random([]Value) Value {
	return Number()
}
