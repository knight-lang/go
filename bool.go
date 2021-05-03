package knight

import (
	"fmt"
	"strconv"
)

type Boolean bool

func (b Boolean) Run() (Value, error) {
	return b, nil
}

func (b Boolean) Dump() {
	fmt.Printf("Boolean(%s)", b)
}

func (b Boolean) Bool() bool {
	return bool(b)
}

func (b Boolean) Int() int {
	if b {
		return 1
	}

	return 0
}

func (b Boolean) String() string {
	return strconv.FormatBool(bool(b))
}
