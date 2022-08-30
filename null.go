package knight

import (
	"fmt"
)

type Null struct{}

var _ Literal = Null{}

func (n Null) Run() (Value, error) { return n, nil }
func (_ Null) Dump()               { fmt.Print("Null()") }
func (_ Null) ToBoolean() Boolean  { return false }
func (_ Null) ToNumber() Number    { return 0 }
func (_ Null) ToText() Text        { return "null" }
func (_ Null) ToList() List        { return nil }
