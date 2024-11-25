package knight

import (
	"fmt"
	"unicode/utf8"
	"unicode"
	"strings"
)

// String is the string type within Knight.
//
// Knight's required encoding is a subset of ASCII; since `string` is UTF-8 encoded, it's already a
// superset, and thus is compliant.
type String string

// Compile-time assertion that `String`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = String("")
var _ Value = String("")

// Run simply returns the text unchanged.
func (s String) Run() (Value, error) {
	return s, nil
}

// Dump prints a debugging representation of `s` to stdout.
func (s String) Dump() {
	fmt.Printf("%q", s)
}

// ToBoolean returns whether `s` is nonempty.
func (s String) ToBoolean() Boolean {
	return s != ""
}

// ToInteger converts `s` to an integer as defined by the knight spec.
func (s String) ToInteger() Integer {
	var ret Integer
	fmt.Sscanf(strings.TrimLeftFunc(string(s), unicode.IsSpace), "%d", &ret)
	return ret
}

// ToString simply returns `s`.
func (s String) ToString() String {
	return s
}

// FirstRune returns the first rune of `s` along with a `String` with that rune removed.
//
// If `s` is empty, this function panics.
func (s String) FirstRune() (rune, String) {
	if s == "" {
		panic("FirstRune called on an empty string")
	}

	r, i := utf8.DecodeRuneInString(string(s))
	return r, s[i:]
}

// ToList returns a `List` of all the `rune`s within `s`.
func (s String) ToList() List {
	list := make(List, 0, utf8.RuneCountInString(string(s)))

	for s != "" {
		var rune rune
		rune, s = s.FirstRune()
		list = append(list, String(rune))
	}

	return list
}
