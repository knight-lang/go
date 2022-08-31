package knight

import (
	"fmt"
)

// Variable represents a variable within Knight.
type Variable struct {
	name  string
	value Value
}

// Compile-time assertion that `Variable`s implements the `Value` interface.
var _ Value = &Variable{}

// Run looks up the last-assigned value for `v`, returning an error if `v` is unassigned.
func (v *Variable) Run() (Value, error) {
	if v.value == nil {
		return nil, fmt.Errorf("undefined variable %q encountered.", v.name)
	}

	return v.value, nil
}

// Dump prints a debug representation of `v` to stdout.
func (v *Variable) Dump() {
	fmt.Printf("Variable(%s)", v.name)
}

// Assign replaces the old value for `v` with `value`.
func (v *Variable) Assign(value Value) {
	v.value = value
}
