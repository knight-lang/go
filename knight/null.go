package knight

import (
	"fmt"
)

// Null is the null type within knight, and is simply an empty struct.
type Null struct{}

// Compile-time assertion that `Null` implements the `Convertible` and `Value` interfaces.
var _ Convertible = Null{}
var _ Value = Null{}

// Run simply returns `n` unchanged.
func (n Null) Run() (Value, error) {
	return n, nil
}

// Dump simply prints `Null()` to stdout.
func (_ Null) Dump() {
	fmt.Print("Null()")
}

// ToBoolean simply returns false.
func (_ Null) ToBoolean() Boolean {
	return false
}

// ToNumber simply returns zero.
func (_ Null) ToNumber() Number {
	return 0
}

// ToText simply returns `"null"`.
func (_ Null) ToText() Text {
	return "null"
}

// ToList simply returns an empty list.
func (_ Null) ToList() List {
	return nil
}
