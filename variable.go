package knight

import "fmt"

type Variable struct {
	name  string
	value Value
}

var variables map[string]*Variable = make(map[string]*Variable)

func NewVariable(name string) *Variable {
	val, ok := variables[name]

	if !ok {
		val = &Variable{name: name}
		variables[name] = val
	}

	return val
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
