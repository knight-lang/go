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

// Dump simply prints "null" to stdout.
func (_ Null) Dump() {
	fmt.Print("null")
}

// Execute simply returns the null unchanged.
func (n Null) Execute() (Value, error) {
	return n, nil
}

// ToBool simply returns false.
func (_ Null) ToBool() (bool, error) {
	return false, nil
}

// ToInt simply returns 0.
func (_ Null) ToInt() (int, error) {
	return 0, nil
}

// ToString simply returns an empty string.
func (_ Null) ToString() (string, error) {
	return "", nil
}

// ToSlice simply returns an empty list.
func (_ Null) ToSlice() ([]Value, error) {
	return nil, nil
}
