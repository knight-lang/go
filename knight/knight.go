package knight

import (
	"fmt"
)

// NothingToParse is the error that's returned when `Play` is given a source string that's empty or
// just comments and whitespace.
var NothingToParse = fmt.Errorf("nothing to parse")

// Play parses `source` as Knight code, and then executes it.
//
// Note that each call to `Play` uses a different set of variables. USe `PlayWithEnvironment` if
// you want to reuse variables.
func Play(source string) (Value, error) {
	env := NewEnvironment()
	return PlayWithEnvironment(source, &env)
}

// PlayWithEnvironment parses `source` as Knight code, and then executes it.
func PlayWithEnvironment(source string, env *Environment) (Value, error) {
	parser := NewParser(source)
	value, err := parser.Parse(env)
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, NothingToParse
	}

	return value.Run()
}
