package knight

import (
	"fmt"
	"strconv"
)

// Integer is the numeric type within Knight, and is simply a wrapper around `int64`.
type Integer int64

// Compile-time assertion that `Integer`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = Integer(0)
var _ Value = Integer(0)

// Run simply returns `n` unchanged.
func (n Integer) Run() (Value, error) {
	return n, nil
}

// Dump prints a debugging representation of `n` to stdout.
func (n Integer) Dump() {
	fmt.Printf("%d", n)
}

// ToBoolean returns whether `n` is nonzero.
func (n Integer) ToBoolean() Boolean {
	return n != 0
}

// ToInteger simply returns `n` unchanged.
func (n Integer) ToInteger() Integer {
	return n
}

// ToText returns the string representation of `n`.
func (n Integer) ToText() Text {
	return Text(strconv.Itoa(int(n)))
}

// ToList returns the the list of digits for `n`.
//
// While not required by the specs, if `n` is negative, `-1` is prepended to the list of digits.
func (n Integer) ToList() List {
	if n == 0 {
		return List{n}
	}

	// TODO: maybe this could be optimized?
	var list List
	for n != 0 {
		list = append(List{n % 10}, list...)
		n /= 10
	}

	return list
}
