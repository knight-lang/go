package knight

import (
	"fmt"
)

// Boolean is the boolean type within Knight.
//
// Knight code can access `Boolean`s via the functions which return booleans (notably `TRUE` and
// `FALSE`, along with comparison operators).
type Boolean bool

// Compile-time assertion that `Boolean` implements the `Convertible` and `Value` interfaces.
var _ Convertible = Boolean(true)
var _ Value = Boolean(true)

// Run simply returns a boolean unchanged.
func (b Boolean) Run() (Value, error) {
	return b, nil
}

// Dump prints a debugging representation of the boolean to stdout.
func (b Boolean) Dump() {
	fmt.Print(b.ToText())
}

// ToBoolean simply returns the boolean unchanged.
func (b Boolean) ToBoolean() Boolean {
	return b
}

// ToInteger returns `1` if the boolean is true, `0` otherwise.
func (b Boolean) ToInteger() Integer {
	if b {
		return 1
	}

	return 0
}

// ToText returns the string representation of the boolean.
func (b Boolean) ToText() Text {
	if b {
		return "true"
	}
	return "false"
}

// ToList returns an empty `List` when the boolean is false, or a list of just the boolean when
// the boolean is true.
func (b Boolean) ToList() List {
	if b {
		return List{b}
	}

	return nil
}
