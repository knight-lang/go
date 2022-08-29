package knight

import (
	"fmt"
	"strconv"
)

type Number int64

func (n Number) Run() (Value, error) {
	return n, nil
}

func (n Number) Dump() {
	fmt.Printf("Number(%d)", n)
}

func (n Number) Bool() bool     { return n != 0 }
func (n Number) Int() int       { return int(n) }
func (n Number) String() string { return strconv.Itoa(int(n)) }
func (n Number) List() []Value {
	if n < 0 {
		panic("negative value given to list")
	}

	if n == 0 {
		return []Value{n}
	}

	// TODO: maybe this could be optimized?
	var list []Value
	for n != 0 {
		list = append([]Value{n % 10}, list...)
		n /= 10
	}

	return list
}
