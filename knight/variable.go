package knight

import (
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
	name  string // the name of the variable; never changed after the Variable is created.
	value Value  // the current value of the variable. a `nil` value indicates the Variable is unset.
}

// Compile-time assertion that Variable implements the Value interface.
var _ Value = &Variable{}

// variablesMap is the a global variable that's a list of all known variables. It's used by the
// NewVariable function to ensure that all variables of the same name point to the same Variable.
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

// Dump prints a debug representation of the variable to stdout.
func (v *Variable) Dump() {
	fmt.Printf("Variable(%s)", v.name)
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

// Assign replaces the old value for the variable with the new value. Panics if value is nil.
func (v *Variable) Assign(value Value) {
	if value == nil {
		panic("[BUG] Variable.Assign called with a nil value?")
	}

	v.value = value
}

// ToString converts the Variable's last assigned value to a String, returning an error if the
// variable hasn't been assigned yet.
func (v *Variable) ToString() (String, error) {
	ran, err := v.Run()
	if err != nil {
		return "", err
	}

	return ran.ToString()
}

// ToInteger converts the Variable's last assigned value to an Integer, returning an error if the
// variable hasn't been assigned yet.
func (v *Variable) ToInteger() (Integer, error) {
	ran, err := v.Run()
	if err != nil {
		return 0, err
	}

	return ran.ToInteger()
}

// ToBoolean converts the Variable's last assigned value to a Boolean, returning an error if the
// variable hasn't been assigned yet.
func (v *Variable) ToBoolean() (Boolean, error) {
	ran, err := v.Run()
	if err != nil {
		return false, err
	}

	return ran.ToBoolean()
}

// ToList converts the Variable's last assigned value to a List, returning an error if the
// variable hasn't been assigned yet.
func (v *Variable) ToList() (List, error) {
	ran, err := v.Run()
	if err != nil {
		return nil, err
	}

	return ran.ToList()
}
