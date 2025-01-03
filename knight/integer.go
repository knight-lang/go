package knight

import (
	"fmt"
	"strconv"
)

// Integer is the numeric type within Knight.
//
// The Knight specs only require support for `int32`, but this implementation supports `int` (which
// is `int64` on 64-bit platforms, but `int32` on 32-bit platforms) as an extension.
type Integer int

// Compile-time assertion that Integer implements the Value interface.
var _ Value = Integer(0)

// Dump prints the integer in base-10 to stdout.
func (i Integer) Dump() {
	fmt.Printf("%d", i)
}

// Execute simply returns the integer unchanged.
func (i Integer) Execute() (Value, error) {
	return i, nil
}

// ToBool returns whether the integer is nonzero.
func (i Integer) ToBool() (bool, error) {
	return i != 0, nil
}

// ToInt simply returns the integer unchanged.
func (i Integer) ToInt() (int, error) {
	return int(i), nil
}

// ToString returns the string representation of the integer in base-10.
func (i Integer) ToString() (string, error) {
	return strconv.Itoa(int(i)), nil
}

// ToSlice returns the digits of the integer in base-10 format.
func (i Integer) ToSlice() ([]Value, error) {
	// Special case for when we're just given 0
	if i == 0 {
		return List{i}, nil
	}

	// Knight 3.0 says that negative integers -> list is undefined behaviour. As an extension, this
	// implementation supports this conversion. (And, with no extra cost too; the algorithm that's
	// used below to get the digit just happens to work on negative numbers too.) It'd also be
	// totally valid to just return an error indicating the conversion isn't supported.

	var list List
	for i != 0 {
		list = append(List{i % 10}, list...)
		i /= 10
	}

	return list, nil
}
