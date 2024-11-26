package knight

import (
	"fmt"
)

// Boolean is the boolean type within Knight.
//
// Knight code can access `Boolean`s via the functions which return booleans (notably `TRUE` and
// `FALSE`, along with comparison operators).
type Boolean bool

// Compile-time assertion that `Boolean` implements the `Value` interface.
var _ Value = Boolean(true)

// Run simply returns a boolean unchanged.
func (b Boolean) Run() (Value, error) {
	return b, nil
}

// Dump prints a debugging representation of the boolean to stdout.
func (b Boolean) Dump() {
	string, _ := b.ToString()
	fmt.Print(string)
}

// ToBoolean simply returns the boolean unchanged.
func (b Boolean) ToBoolean() (Boolean, error) {
	return b, nil
}

// ToInteger returns `1` if the boolean is true, `0` otherwise.
func (b Boolean) ToInteger() (Integer, error) {
	if b {
		return 1, nil
	}
	return 0, nil
}

// ToString returns the string representation of the boolean.
func (b Boolean) ToString() (String, error) {
	if b {
		return "true", nil
	}
	return "false", nil
}

// ToList returns an empty `List` when the boolean is false, or a list of just the boolean when
// the boolean is true.
func (b Boolean) ToList() (List, error) {
	if b {
		return List{b}, nil
	}
	return nil, nil
}
