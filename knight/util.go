package knight

import (
	"fmt"
)

func bug(format string, args ...any) {
	// if len(args) != fun.arity {
	panic(fmt.Sprintf("[BUG] " + format, args...))
}
