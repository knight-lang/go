package knight

import "fmt"

type Number int

func (n Number) Run() (Value, error) {
	return n, nil
}

func (n Number) Dump() {
	fmt.Printf("Number(%d)", n)
}

func (n Number) Bool() bool {
	return n != 0
}

func (n Number) Int() int {
	return int(n)
}

func (n Number) String() string {
	return string(n)
}
