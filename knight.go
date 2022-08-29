package knight

import (
	"fmt"
)

func Run(input string, e *Environment) (Value, error) {
	parser := NewParser(input)
	value, err := parser.Parse(e)

	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, fmt.Errorf("nothing to parse")
	}

	return value.Run()
}
