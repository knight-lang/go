package knight

import (
	"fmt"
)

// Play parses `source` as Knight code, and then executes it.
//
// Note that each call to `Play` uses a different set of variables. USe `PlayWithEnvironment` if
// you want to reuse variables.
func Play(source string) (Value, error) {
	return PlayWithEnvironment(source, NewEnvironment())
}

// PlayWithEnvironment parses `source` as Knight code, and then executes it.
func PlayWithEnvironment(source string, env *Environment) (Value, error) {
	parser := NewParser(source)
	value, err := parser.Parse(env)
	if err != nil {
		return nil, fmt.Errorf("compile error: %v", err)
	}

	result, err := value.Run()
	if err != nil {
		return nil, fmt.Errorf("runtime error: %v", err)
	}

	return result, nil
}
