package knight

import (
	"fmt"
)

type Boolean bool

var _ Literal = Boolean(true)

func (b Boolean) Run() (Value, error) {
	return b, nil
}

func (b Boolean) Dump() {
	fmt.Printf("Boolean(%s)", b)
}

func (b Boolean) Bool() bool {
	return b
}

func (b Boolean) Int() int {
	if b {
		return 1
	}
	return 0
}

func (b Boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (b Boolean) List() List {
	if b {
		return List{b}
	}
	return nil
}
