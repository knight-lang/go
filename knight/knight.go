package knight

import (
	"fmt"
)

// shouldSupportKnightVersion_2_0_1 is a toggle for whether Knight 2.0.1 should be supported. There
// are some breaking changes in 2.0.1 for Knight code (however, not for implementations---all 2.0.1
// implementations are valid Knight 2.1 implementations). If this is set to true, those breaking
// changes aren't supported.
const shouldSupportKnightVersion_2_0_1 = true

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
