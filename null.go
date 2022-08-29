package knight

import (
	"fmt"
)

type Null struct{}

var _ Literal = Null{}

func (n Null) Run() (Value, error) { return n, nil }
func (n Null) Dump()               { fmt.Print("Null()") }
func (_ Null) Bool() bool          { return false }
func (_ Null) Int() int            { return 0 }
func (_ Null) String() string      { return "null" }
func (_ Null) List() List          { return nil }
