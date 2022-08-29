package knight

import (
	"fmt"
)

type List []Value

func (l List) Run() (Value, error) {
	return l, nil
}

func (l List) Dump() {
	fmt.Printf("List(%q)", l)
}

func (l List) Bool() bool     { return len(l) != 0 }
func (l List) Int() int       { return len(l) }
func (l List) String() string { panic("todo") }
func (l List) List() []Value  { return []Value(l) }

// func (l List) Append(rhs []Value) List {
// 	return List(append([]Value(l), rhs...))
// }

// func (l List) Repeat(rhs []Value) List {
// 	return List(append([]Value(l), rhs...))
// }
