package knight

import (
	"fmt"
	"strings"
)

// List is the list type within Knight
//
// It's actually just a wrapper around `[]Value`.
type List []Value

// Compile-time assertion that `List`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = List{}
var _ Value = List{}

// Run simply returns the list unchanged.
func (l List) Run() (Value, error) {
	return l, nil
}

// Dump prints a debugging representation of the list to stdout.
func (l List) Dump() {
	fmt.Print("[")

	for i, element := range l {
		// Don't print a comma for the first argument
		if i != 0 {
			fmt.Print(", ")
		}

		element.Dump()
	}

	fmt.Print("]")
}

// ToBoolean returns whether or not the list is empty.
func (l List) ToBoolean() Boolean {
	return len(l) != 0
}

// ToInteger returns the list's length length.
func (l List) ToInteger() Integer {
	return Integer(len(l))
}

// ToText returns the list converted to a string by adding a newline between each element.
func (l List) ToText() Text {
	return Text(l.Join("\n"))
}

// ToList simply returns the list unchanged.
func (l List) ToList() List {
	return l
}

// Join concatenates all the elements of the list together into a big string, with `separator`
// interspersed between the elements.
func (l List) Join(separator string) string {
	// Use a `strings.Builder` for efficiency, as we'll be doing multiple concatenations.
	var sb strings.Builder

	for i, element := range l {
		// Don't add the separator during the first iteration
		if i != 0 {
			sb.WriteString(separator)
		}

		// Add the element to the end. Note that `element.(Convertible)` will panic if `element`
		// doesn't implement `Convertible`, which can happen if the element's the return value of
		// `BLOCK`. However, this is OK, as the Knight specs don't require `BLOCK`'s return values to
		// be convertible to strings.
		sb.WriteString(string(element.(Convertible).ToText()))
	}

	return sb.String()
}
