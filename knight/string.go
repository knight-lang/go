package knight

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// String is the string type within Knight.
//
// Knight's required encoding is a subset of ASCII; since `string` is UTF-8 encoded, it's already a
// superset, and thus is compliant.
type String string

// Compile-time assertion that `String`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = String("")
var _ Value = String("")

// Run simply returns the string unchanged.
func (s String) Run() (Value, error) {
	return s, nil
}

// Dump prints a debugging representation of the string to stdout.
func (s String) Dump() {
	// It just so happens that golang's `%q` specifier exactly matches what Knight's `DUMP` expects.
	fmt.Printf("%q", s)
}

// ToBoolean returns whether the string is nonempty.
func (s String) ToBoolean() Boolean {
	return s != ""
}

// ToInteger converts the string to an integer as defined by the knight spec.
func (s String) ToInteger() Integer {
	var ret Integer
	fmt.Sscanf(strings.TrimLeftFunc(string(s), unicode.IsSpace), "%d", &ret)
	return ret
}

// ToString simply returns the string unchanged.
func (s String) ToString() String {
	return s
}

// StringIsEmpty is an error that's returned by `SplitFirstRune` when a string is empty.
var StringIsEmpty = errors.New("SplitFirstRune called on an empty string")

// SplitFirstRune returns the first rune of the string and a `String` with that rune removed.
//
// If the string is empty, this function returns an `error`.
func (s String) SplitFirstRune() (rune, String, error) {
	if s == "" {
		return 0, "", StringIsEmpty
	}

	rune, idx := utf8.DecodeRuneInString(string(s))
	return rune, s[idx:], nil
}

// ToList returns a list of all the `rune`s within the string.
func (s String) ToList() List {
	list := make(List, 0, utf8.RuneCountInString(string(s)))

	for s != "" {
		// We know that `SplitFirstRune` can't fail as we just checked to see if it was empty.
		var rune rune
		rune, s, _ = s.SplitFirstRune()

		list = append(list, String(rune))
	}

	return list
}
