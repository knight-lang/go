package knight

import (
	"fmt"
)

type Variable struct {
	name  string
	value Value
}

type Environment struct {
	variables map[string]*Variable
}

func NewEnvironment() Environment {
	return Environment{
		variables: make(map[string]*Variable),
	}
}

func (e *Environment) Lookup(name string) *Variable {
	if variable, ok := e.variables[name]; ok {
		return variable
	}

	variable := &Variable{name: name}
	e.variables[name] = variable
	return variable
}

func (v *Variable) Run() (Value, error) {
	if v.value == nil {
		return nil, fmt.Errorf("undefined variable %q encountered.", v.name)
	}

	return v.value, nil
}

func (v *Variable) Dump() {
	fmt.Printf("Variable(%s)", v.name)
}

func (v *Variable) Assign(value Value) {
	v.value = value
}
