package knight

import (
	"fmt"
	"strconv"
)

// Integer is the numeric type within Knight.
//
// Technically the specs only require supporting up to `int32`. However, we support up to `int64`,
// as a convenience to end-users.
type Integer int64

// Compile-time assertion that `Integer`s implements the `Convertible` and `Value` interfaces.
var _ Convertible = Integer(0)
var _ Value = Integer(0)

// Run simply returns the integer unchanged.
func (i Integer) Run() (Value, error) {
	return i, nil
}

// Dump prints a debugging representation of `i` to stdout.
func (i Integer) Dump() {
	fmt.Printf("%d", i)
}

// ToBoolean returns whether the integer is nonzero.
func (i Integer) ToBoolean() Boolean {
	return i != 0
}

// ToInteger simply returns the integer unchanged.
func (i Integer) ToInteger() Integer {
	return i
}

// ToString returns the string representation of the integer.
func (i Integer) ToString() String {
	return String(strconv.FormatInt(int64(i), 10))
}

// ToList returns the digits of the integer in base-10 format.
//
// While not required by the specs, if the integer is negative, each digit is negated. (i.e. 
// `Integer(-123).ToList()` is `{-1, -2, -3}`).
func (i Integer) ToList() List {
	// Special case for when we're just given 0
	if i == 0 {
		return List{i}
	}

	var list List
	for i != 0 {
		list = append(List{i % 10}, list...)
		i /= 10
	}

	return list
}
