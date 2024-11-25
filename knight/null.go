package knight

import (
	"fmt"
)

// Null is the null type within knight.
//
// It's an empty struct, as there's only one null instance, and it doesn't require additional info.
//
// Knight code can access `Null` via the `NULL` function, and a handful of other functions which
// return `NULL` (eg `OUTPUT`.)
type Null struct{}

// Compile-time assertion that `Null` implements the `Convertible` and `Value` interfaces.
var _ Convertible = Null{}
var _ Value = Null{}

// Run simply returns the null unchanged.
func (n Null) Run() (Value, error) {
	return n, nil
}

// Dump simply prints `null` to stdout.
func (_ Null) Dump() {
	fmt.Print("null")
}

// ToBoolean simply returns false.
func (_ Null) ToBoolean() Boolean {
	return false
}

// ToInteger simply returns zero.
func (_ Null) ToInteger() Integer {
	return 0
}

// ToString simply returns an empty string.
func (_ Null) ToString() String {
	return ""
}

// ToList simply returns an empty list.
func (_ Null) ToList() List {
	return nil
}
