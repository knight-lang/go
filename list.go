package knight

import (
	"fmt"
	"strings"
)

type List []Value

var _ Literal = List{}

func (l List) Run() (Value, error) {
	return l, nil
}

func (l List) Dump() {
	fmt.Printf("List(%q)", l)
}

func (l List) Bool() bool {
	return len(l) != 0
}

func (l List) Int() int       { return len(l) }
func (l List) String() string { return l.Join("\n") }
func (l List) List() List     { return l }

func (l List) Join(sep string) string {
	var sb strings.Builder

	for i, ele := range l {
		if i != 0 {
			sb.WriteString(sep)
		}

		sb.WriteString(ele.(interface{ String() string }).String())
	}

	return sb.String()
}
