package knight

import (
	"fmt"
)

// Variable represents a variable within Knight.
//
// Normally this type is opaque to the end-user, as whenever it's executed it gives a value (or an
// error) back. However, this implementation is designed such that `BLOCK <expr>` returns `<expr>`
// unchanged, so it's possible for Knight programs to "interact" with variables via `DUMP BLOCK hi`.
// 
// Knight's specs are very explicit about what `BLOCK`'s return value can be used for---not much,
// other than `CALL`. As such, the `Variable` type doesn't need to implement `Convertible`.
//
// Variables are exclusively created via the `Environment` type's `Lookup` function. See it for more
// details.
type Variable struct {
	// name of the variable.
	name  string

	// value is the variable's current value
	//
	// When we initialize a new `Variable` (in `NewVariable`), this field is set to the default value
	// of `nil`. Since `nil` isn't a valid `Value`, we can use it to see if a program's accessing an
	// undefined value. (Technically this isn't required, because the Knight spec says that accessing
	// unassigned variables is undefined behaviour, but this way we can provide nice error messages.)
	value Value
}

// Compile-time assertion that `Variable`s implements the `Value` interface.
var _ Value = &Variable{}

// variablesMap is the list of all registered variables. It's used by `NewVariable` to ensure that
// all variables of the same name point to the same value.
var variablesMap map[string]*Variable

// NewVariable returns the variable corresponding to `name`; if it doesn't exist, it's created.
func NewVariable(name string) *Variable {
	// If the variable already exists, then return it.
	if variable, ok := variablesMap[name]; ok {
		return variable
	}

	// The variable doesn't exist. Create it, add it to `variablesMap`, then return it.
	variable := &Variable{ name: name }
	variablesMap[name] = variable
	return variable
}


// Run looks up the last-assigned value for the variable, returning an error if the variable hasn't
// been assigned yet.
func (v *Variable) Run() (Value, error) {
	if v.value == nil {
		return nil, fmt.Errorf("undefined variable %q encountered.", v.name)
	}

	return v.value, nil
}

// Dump prints a debug representation of the variable to stdout.
func (v *Variable) Dump() {
	fmt.Printf("Variable(%s)", v.name)
}

// Assign replaces the old value for the variable with the new `value`.
func (v *Variable) Assign(value Value) {
	v.value = value
}
