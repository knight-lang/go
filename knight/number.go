package knight

import (
	"fmt"
	"strconv"
)

// Number is the numeric type within Knight, and is simply a wrapper around `int64`.
type Number int64

// Compile-time assertion that `Number`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = Number(0)
var _ Value = Number(0)

// Run simply returns `n` unchanged.
func (n Number) Run() (Value, error) {
	return n, nil
}

// Dump prints a debugging representation of `n` to stdout.
func (n Number) Dump() {
	fmt.Printf("Number(%d)", n)
}

// ToBoolean returns whether `n` is nonzero.
func (n Number) ToBoolean() Boolean {
	return n != 0
}

// ToNumber simply returns `n` unchanged.
func (n Number) ToNumber() Number {
	return n
}

// ToText returns the string representation of `n`.
func (n Number) ToText() Text {
	return Text(strconv.Itoa(int(n)))
}

// ToList returns the the list of digits for `n`.
//
// While not required by the specs, if `n` is negative, `-1` is prepended to the list of digits.
func (n Number) ToList() List {
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
