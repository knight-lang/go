package knight

import (
	"fmt"
)

// Null is the null type within knight.
//
// Knight code can access Null via the `NULL` function, and a handful of other functions which
// return `NULL` (eg `OUTPUT`).
type Null struct{} // empty struct because there's only ever one null

// Compile-time assertion that Null implements the Value interface.
var _ Value = Null{}

// Run simply returns the null unchanged.
func (n Null) Run() (Value, error) {
	return n, nil
}

// Dump simply prints "null" to stdout.
func (_ Null) Dump() {
	fmt.Print("null")
}

// ToBoolean simply returns false.
func (_ Null) ToBoolean() (Boolean, error) {
	return false, nil
}

// ToInteger simply returns 0.
func (_ Null) ToInteger() (Integer, error) {
	return 0, nil
}

// ToString simply returns an empty string.
func (_ Null) ToString() (String, error) {
	return "", nil
}

// ToList simply returns an empty list.
func (_ Null) ToList() (List, error) {
	return nil, nil // nil is equivalent to an empty slice (and Lists are wrappers around slices).
}
