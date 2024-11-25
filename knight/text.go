package knight

import (
	"fmt"
	"unicode/utf8"
	"unicode"
	"strings"
)

// Text is the string type within Knight, and is simply a wrapper around `string`.
type Text string

// Compile-time assertion that `Text`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = Text("")
var _ Value = Text("")

// Run simply returns the `t` unchanged.
func (t Text) Run() (Value, error) {
	return t, nil
}

// Dump prints a debugging representation of `t` to stdout.
func (t Text) Dump() {
	fmt.Printf("%q", t)
}

// ToBoolean returns whether `t` is nonempty.
func (t Text) ToBoolean() Boolean {
	return t != ""
}

// ToInteger converts `t` to an integer as defined by the knight spec.
func (t Text) ToInteger() Number {
	var ret Number
	fmt.Sscanf(strings.TrimLeftFunc(string(t), unicode.IsSpace), "%d", &ret)
	return ret
}

// ToText simply returns `t`.
func (t Text) ToText() Text {
	return t
}

// FirstRune returns the first rune of `t` along with a `Text` with that rune removed.
//
// If `t` is empty, this function panics.
func (t Text) FirstRune() (rune, Text) {
	if t == "" {
		panic("FirstRune called on an empty string")
	}

	r, i := utf8.DecodeRuneInString(string(t))
	return r, t[i:]
}

// ToList returns a `List` of all the `rune`s within `t`.
func (t Text) ToList() List {
	list := make(List, 0, utf8.RuneCountInString(string(t)))

	for t != "" {
		var rune rune
		rune, t = t.FirstRune()
		list = append(list, Text(rune))
	}

	return list
}
