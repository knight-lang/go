package knight

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// String is the type for holding text within Knight.
//
// Knight's specs only require implementations to support a specific subset of ASCII. However, as a
// convenience to end-users, String *also* supports all of UTF-8.
type String string

// Compile-time assertion that String implements the Value interface.
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
func (s String) ToBoolean() (Boolean, error) {
	return s != "", nil
}

// ToInteger converts the string to an integer as defined by the knight spec.
//
// More specifically, this is equivalent to matching the string against the regex `/^\s+([-+]?\d+)/`
// and converting the first capture group (the `[-+]?\d+`) to a string. If the regex doesn't match,
// then zero is used.
func (s String) ToInteger() (Integer, error) {
	var ret Integer
	fmt.Sscanf(strings.TrimLeftFunc(string(s), unicode.IsSpace), "%d", &ret)
	return ret, nil
}

// ToString simply returns the string unchanged.
func (s String) ToString() (String, error) {
	return s, nil
}

// StringIsEmpty is an error that's returned by SplitFirstRune when a string is empty.
var StringIsEmpty = errors.New("SplitFirstRune called on an empty string")

// SplitFirstRune returns the first rune of the string and a String with that rune removed.
//
// If the string is empty, this function returns an error.
func (s String) SplitFirstRune() (rune, String, error) {
	if s == "" {
		return 0, "", StringIsEmpty
	}

	rune, idx := utf8.DecodeRuneInString(string(s))
	return rune, s[idx:], nil
}

// ToList returns a list of all the runes within the string.
func (s String) ToList() (List, error) {
	list := make(List, 0, utf8.RuneCountInString(string(s)))

	for s != "" {
		// We know that SplitFirstRune can't fail as we just checked to see if it was empty.
		var rune rune
		rune, s, _ = s.SplitFirstRune()

		list = append(list, String(rune))
	}

	return list, nil
}
