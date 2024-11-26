package knight

import (
	"fmt"
	"strings"
)

// List is the list type within Knight
//
// It's actually just a wrapper around `[]Value`.
type List []Value

// Compile-time assertion that `List`s implements the `Value` interfaces.
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
func (l List) ToBoolean() (Boolean, error) {
	return len(l) != 0, nil
}

// ToInteger returns the list's length length.
func (l List) ToInteger() (Integer, error) {
	return Integer(len(l)), nil
}

// ToString returns the list converted to a string by adding a newline between each element.
func (l List) ToString() (String, error) {
	joined, err := l.Join("\n")
	if err != nil {
		return "", err
	}

	return String(joined), nil
}

// ToList simply returns the list unchanged.
func (l List) ToList() (List, error) {
	return l, nil
}

// Join concatenates all the elements of the list together into a big string, with `separator`
// interspersed between the elements.
func (l List) Join(separator string) (string, error) {
	// Use a `strings.Builder` for efficiency, as we'll be doing multiple concatenations.
	var sb strings.Builder

	for i, element := range l {
		// Don't add the separator during the first iteration
		if i != 0 {
			sb.WriteString(separator)
		}

		repr, err := element.ToString()
		if err != nil {
			return "", err
		}

		sb.WriteString(string(repr))
	}

	return sb.String(), nil
}
