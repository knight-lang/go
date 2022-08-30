package knight

import (
	"fmt"
	"unicode/utf8"
)

// Text is the string type within Knight, and is simply a wrapper around `string`.
type Text string

// Compile-time assertion that `Text`s implement the `Literal` interface.
var _ Literal = Text("")

// Run simply returns the `t` unchanged.
func (t Text) Run() (Value, error) {
	return t, nil
}

// Dump prints a debugging representation of `t` to stdout.
func (t Text) Dump() {
	fmt.Printf("String(%s)", t)
}

// ToBoolean returns whether `t` is nonempty.
func (t Text) ToBoolean() Boolean {
	return t != ""
}

// Int converts `t` to an integer as defined by the knight spec.
func (t Text) ToNumber() Number {
	var ret Number
	fmt.Sscanf(string(t), "%d", &ret)
	return ret
}

// String simply unwraps `t`.
func (t Text) ToText() Text {
	return t
}

func (t Text) FirstRune() (rune, Text) {
	r, i := utf8.DecodeRuneInString(string(t))
	return r, t[i:]
}

// List returns a `List` of all the runes within `t`.
func (t Text) ToList() List {
	list := make(List, 0, utf8.RuneCountInString(string(t)))

	for t != "" {
		var rune rune
		rune, t = t.FirstRune()
		list = append(list, Text(rune))
	}

	return list
}
