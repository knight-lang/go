package knight

import (
	"fmt"
)

// Evaluate parses source as Knight code, and then executes it. Any errors that occur when parsing
// or executing the code are returned.
func Evaluate(source string) (Value, error) {
	parser := NewParser(source)

	value, err := parser.ParseNextValue()
	if err != nil {
		return nil, fmt.Errorf("compile error: %v", err)
	}

	result, err := value.Execute()
	if err != nil {
		return nil, fmt.Errorf("runtime error: %v", err)
	}

	return result, nil
}
