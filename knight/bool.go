package knight

import (
	"fmt"
)

// Boolean is the boolean type within Knight, and is simply a wrapper around `bool`.
type Boolean bool

// Compile-time assertion that `Boolean`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = Boolean(true)
var _ Value = Boolean(true)

// Run simply returns `b` unchanged.
func (b Boolean) Run() (Value, error) {
	return b, nil
}

// Dump prints a debugging representation of `b` to stdout.
func (b Boolean) Dump() {
	fmt.Printf("Boolean(%s)", b.ToText())
}

// ToBoolean simply returns `b` unchanged.
func (b Boolean) ToBoolean() Boolean {
	return b
}

// ToNumber returns `1` if `b` is true, `0` otherwise.
func (b Boolean) ToNumber() Number {
	if b {
		return 1
	}
	return 0
}

// ToText returns the string representation of `b`.
func (b Boolean) ToText() Text {
	if b {
		return "true"
	}
	return "false"
}

// ToList returns `b`'s `List` representation: An empty `List` when false and `List{b}` when true.
func (b Boolean) ToList() List {
	if b {
		return List{b}
	}
	return nil
}
