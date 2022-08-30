package knight

import (
	"fmt"
	"strings"
)

// List is the list type within Knight, and is simply a wrapper around `[]Value`.
type List []Value

// Compile-time assertion that `List`s implement the `Literal` interface.
var _ Literal = List{}

// Run simply returns `l` unchanged.
func (l List) Run() (Value, error) {
	return l, nil
}

// Dump prints a debugging representation of `l` to stdout.
func (l List) Dump() {
	fmt.Print("List(")

	for i, ele := range l {
		if i != 0 {
			fmt.Print(", ")
		}

		ele.Dump()
	}

	fmt.Print(")")
}

// Bool returns whether `l` is nonempty.
func (l List) ToBoolean() Boolean {
	return len(l) != 0
}

// Int returns `l`'s length.
func (l List) ToNumber() Number {
	return Number(len(l))
}

// String returns `l` converted to a string, with `\n` inserted between each element.
func (l List) ToText() Text {
	return Text(l.Join("\n"))
}

// List simply returns `l`.
func (l List) ToList() List {
	return l
}

func (l List) Join(sep string) string {
	var sb strings.Builder

	for i, ele := range l {
		if i != 0 {
			sb.WriteString(sep)
		}

		sb.WriteString(string(ele.(interface{ ToText() Text }).ToText()))
	}

	return sb.String()
}
