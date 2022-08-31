package knight

// Environment holds all the defined variables for a program.
//
// This isn't needed when executing `Value`s, as there's no way to dynamically look up variables
// within Knight (without extensions).
type Environment struct {
	variables map[string]*Variable
	functions map[rune]*Function
}

// NewEnvironment creates a blank `Environment`.
func NewEnvironment() Environment {
	env := Environment{
		variables: make(map[string]*Variable),
		functions: make(map[rune]*Function),
	}

	populateDefaultFunctions(&env)
	return env
}

// RegisterFunction inserts `fn` into the list of known `Function`s, which are used when parsing.
//
// The previously returned `Function`, if any, is returned.
func (e *Environment) RegisterFunction(fn *Function) *Function {
	old := e.functions[fn.name]
	e.functions[fn.name] = fn

	return old
}

// GetFunction looks up the function associated with `name`, returning `nil` if it doesn't exist.
func (e *Environment) GetFunction(name rune) *Function {
	return e.functions[name]
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
