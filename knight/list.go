package knight

import (
	"fmt"
	"strings"
)

// List is the "container" type within Knight.
//
// Empty lists in Knight are represented via `@`. Lists can be created via `,` (eg `,3`), which
// create a one-element list, or via coercions such as `+ @ 123` (which yields `[1, 2, 3]`.)
type List []Value

// Compile-time assertion that List implements the Value interface.
var _ Value = List{}

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

// Execute simply returns the list unchanged.
func (l List) Execute() (Value, error) {
	return l, nil
}

// ToBoolean returns whether the list is nonempty.
func (l List) ToBoolean() (Boolean, error) {
	return len(l) != 0, nil
}

// ToInteger returns the list's length.
func (l List) ToInteger() (Integer, error) {
	return Integer(len(l)), nil
}

// ToString returns the list converted to a string by adding a newline between each element. This
// will return an error if the list contains elements which aren't convertible to strings, such as
// `BLOCK`'s return value.
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
// interspersed between the elements. An error is returned if an element isn't convertible to a
// string.
func (l List) Join(separator string) (string, error) {
	// Use a `strings.Builder` for efficiency, as we'll be doing multiple concatenations.
	var builder strings.Builder

	for i, element := range l {
		// Don't add the separator during the first iteration
		if i != 0 {
			builder.WriteString(separator)
		}

		// Try to get the string representation, or return an error if there's a problem.
		stringRepresentation, err := element.ToString()
		if err != nil {
			return "", err
		}

		// (We have to cast to `string` as `stringRepresentation` is a `String`, not `string`.)
		builder.WriteString(string(stringRepresentation))
	}

	return builder.String(), nil
}
