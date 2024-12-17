package knight

import (
	"errors"
	"fmt"
)

// Variable represents a variable within Knight code.
//
// Variables are created via the NewVariable function, which ensures that each variable of a given
// name always points to the same underlying Variable struct.
//
// Normally, this type isn't accessible from within Knight programs, as most functions Run their
// arguments before interacting with them. However, the `BLOCK` function has been implemented to
// just return its argument, unevaluated. (So, variables can be accessed via eg `OUTPUT BLOCK foo`.)
// This isn't a problem for spec-compliance, however, as the only valid use for `BLOCK`s are to be
// `CALL`ed, which then will Run the variable anyways.
type Variable struct {
	name  string
	value Value
}

// Compile-time assertion that Variable implements the Value interface.
var _ Value = &Variable{}

// variablesMap is the list of all known variables. It's used by NewVariable to ensure that all
// variables of the same name point to the same Variable.
var variablesMap map[string]*Variable = make(map[string]*Variable)

// NewVariable returns the Variable corresponding to name, creating it if it doesn't exist.
func NewVariable(name string) *Variable {
	// If the variable already exists, then return it.
	if variable, ok := variablesMap[name]; ok {
		return variable
	}

	// The variable doesn't exist. Create it, add it to variablesMap, then return it.
	variable := &Variable{name: name, value: nil}
	variablesMap[name] = variable
	return variable
}

// Run looks up the last-assigned value for the variable, returning an error if the variable hasn't
// been assigned yet.
func (v *Variable) Run() (Value, error) {
	// Assign doesn't allow nil to be assigned to v.value, so we can use nil as a marker for
	// unassigned variables.
	if v.value == nil {
		return nil, fmt.Errorf("undefined variable %q encountered", v.name)
	}

	return v.value, nil
}

// Dump prints a debug representation of the variable to stdout.
func (v *Variable) Dump() {
	fmt.Printf("Variable(%s)", v.name)
}

// Assign replaces the old value for the variable with the new value. Panics if value is nil.
func (v *Variable) Assign(value Value) {
	if value == nil {
		panic("cannot assign nil values")
	}

	v.value = value
}

// Conversions: They always return errors, as variables cannot be converted to other types.
var (
	NoToStringDefinedForVariable  = errors.New("Variable doesn't define string conversions")
	NoToIntegerDefinedForVariable = errors.New("Variable doesn't define integer conversions")
	NoToBooleanDefinedForVariable = errors.New("Variable doesn't define boolean conversions")
	NoToListDefinedForVariable    = errors.New("Variable doesn't define list conversions")
)

func (_ *Variable) ToString() (String, error)   { return "", NoToStringDefinedForVariable }
func (_ *Variable) ToInteger() (Integer, error) { return 0, NoToIntegerDefinedForVariable }
func (_ *Variable) ToBoolean() (Boolean, error) { return false, NoToBooleanDefinedForVariable }
func (_ *Variable) ToList() (List, error)       { return nil, NoToListDefinedForVariable }
