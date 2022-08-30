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

// Environment holds all the defined variables for a program.
//
// This isn't needed when executing `Value`s, as there's no way to dynamically look up variables
// within Knight (without extensions).
type Environment struct {
	variables map[string]*Variable
}

// NewEnvironment creates a blank `Environment`.
func NewEnvironment() Environment {
	return Environment{variables: make(map[string]*Variable)}
}

// Lookup fetches the variable corresponding to `name`. If one doesn't exist, it is created.
func (e *Environment) Lookup(name string) *Variable {
	if variable, ok := e.variables[name]; ok {
		return variable
	}

	variable := &Variable{name: name}
	e.variables[name] = variable
	return variable
}

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
