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

// ToInteger returns `l`'s length.
func (l List) ToInteger() Number {
	return Number(len(l))
}

// ToText returns `l` converted to a string by adding a newline between each element.
func (l List) ToText() Text {
	return Text(l.Join("\n"))
}

// ToList simply returns `l`.
func (l List) ToList() List {
	return l
}

// Join concatenates all the elements of `l` together into a big string, with `sep` interspersed.
func (l List) Join(sep string) string {
	var sb strings.Builder

	for i, ele := range l {
		if i != 0 {
			sb.WriteString(sep)
		}

		sb.WriteString(string(ele.(Convertible).ToText()))
	}

	return sb.String()
}
