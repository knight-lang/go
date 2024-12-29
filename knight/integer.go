package knight

import (
	"fmt"
	"strconv"
)

// Integer is the numeric type within Knight.
//
// The Knight specs only require support for `int32`, but this implementation supports `int64`
// as an extension.
type Integer int64

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

// ToBoolean returns whether the integer is nonzero.
func (i Integer) ToBoolean() (Boolean, error) {
	return i != 0, nil
}

// ToInteger simply returns the integer unchanged.
func (i Integer) ToInteger() (Integer, error) {
	return i, nil
}

// ToString returns the string representation of the integer in base-10.
func (i Integer) ToString() (String, error) {
	return String(strconv.FormatInt(int64(i), 10)), nil
}

// ToList returns the digits of the integer in base-10 format.
func (i Integer) ToList() (List, error) {
	// Special case for when we're just given 0
	if i == 0 {
		return List{i}, nil
	}

	// Knight 3.0 says that negative integers -> list is undefined behaviour. As an extension, this
	// implementation supports this conversion. (And, with no extra cost too; the algorithm that's
	// used below to get the digit just happens to work on negative numbers too.)

	var list List
	for i != 0 {
		list = append(List{i % 10}, list...)
		i /= 10
	}

	return list, nil
}
