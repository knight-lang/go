package knight

import (
	"fmt"
	"strconv"
)

// Integer is the numeric type within Knight.
//
// The Knight specs only require that implementations support up to int32, however we support up
// to int64 as a convenience for end-users.
type Integer int64

// Compile-time assertion that Integer implements the Value interface.
var _ Value = Integer(0)

// Dump prints the integer in base-10 to stdout.
func (i Integer) Dump() {
	fmt.Printf("%d", i)
}

// Run simply returns the integer unchanged.
func (i Integer) Run() (Value, error) {
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

// ToList returns the digits of the integer in base-10 format. If the integer is negative, an error
// is returned.
func (i Integer) ToList() (List, error) {
	// Special case for when we're just given 0
	if i == 0 {
		return List{i}, nil
	}

	// Knight 2.0.1 requires negative integers to be supported, an to have their sign added to each
	// element of the list (ie `-123` would become `[-1, -2, -3]`).
	if !shouldSupportKnightVersion_2_0_1 && i < 0 {
		return nil, fmt.Errorf("attempted to convert a negative integer to list: %d", i)
	}

	var list List
	for i != 0 {
		list = append(List{i % 10}, list...)
		i /= 10
	}

	return list, nil
}
