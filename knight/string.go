package knight

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// String is the type for holding text within Knight.
//
// Knight's specs only require implementations to support a specific subset of ASCII. However, as a
// convenience to end-users, String *also* supports all of Unicode. (We get this for free because we
// use go's `string` type, which supports Unicode.)
type String string

// Compile-time assertion that String implements the Value interface.
var _ Value = String("")

// Dump prints the escaped version of string to stdout.
func (s String) Dump() {
	// It just so happens that golang's `%q` specifier exactly matches what Knight's `DUMP` expects.
	fmt.Printf("%q", s)
}

// Execute simply returns the Stirng unchanged.
func (s String) Execute() (Value, error) {
	return s, nil
}

// ToBool returns whether the string is nonempty.
func (s String) ToBool() (bool, error) {
	return s != "", nil
}

// ToInt converts the string to an integer as defined by the knight spec.
//
// More specifically, this is equivalent to matching the string against the regex `/^\s+([-+]?\d+)/`
// and converting the first capture group (the `[-+]?\d+`) to a string. If the regex doesn't match,
// then zero is used.
func (s String) ToInt() (int, error) {
	// Delete leading whitespace
	trimmed := strings.TrimLeftFunc(string(s), unicode.IsSpace)

	// Parse out the integer. If Scanf fails, parsed stays zero.
	var parsed int
	fmt.Sscanf(trimmed, "%d", &parsed)

	// No errors can occur when converting strings to integers.
	return parsed, nil
}

// ToString simply returns the string unchanged.
func (s String) ToString() (string, error) {
	return string(s), nil
}

// ToSlice returns a list of all the runes within string.
func (s String) ToSlice() ([]Value, error) {
	list := make([]Value, utf8.RuneCountInString(string(s)))

	for idx, rune := range []rune(s) {
		list[idx] = String(rune)
	}

	return list, nil
}
