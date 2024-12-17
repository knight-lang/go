package knight

import (
	"fmt"
)

// Boolean is the boolean type within Knight.
//
// Knight code can access Booleans via functions which return booleans (notably `TRUE` and `FALSE`).
type Boolean bool

// Compile-time assertion that Boolean implements the Value interface.
var _ Value = Boolean(false)

// Dump prints "true" or "false" to stdout.
func (b Boolean) Dump() {
	// In golang, `%t` is for booleans (just like `%d` is for ints and `%s` is for strings)
	fmt.Printf("%t", b)
}

// Execute simply returns the boolean unchanged.
func (b Boolean) Execute() (Value, error) {
	return b, nil
}

// ToBoolean simply returns the boolean unchanged.
func (b Boolean) ToBoolean() (Boolean, error) {
	return b, nil
}

// ToInteger returns 1 if the boolean is true and 0 if it is false.
func (b Boolean) ToInteger() (Integer, error) {
	if b {
		return 1, nil
	}
	return 0, nil
}

// ToString returns "true" or "false"
func (b Boolean) ToString() (String, error) {
	if b {
		return "true", nil
	}
	return "false", nil
}

// ToList returns an empty List when the boolean is false, or a list of just the boolean when
// the boolean is true.
func (b Boolean) ToList() (List, error) {
	if b {
		return List{b}, nil
	}
	return nil, nil
}
