package knight

import (
	"fmt"
	"strconv"
)

type Number int64

var _ Literal = Number(0)

func (n Number) Run() (Value, error) {
	return n, nil
}

func (n Number) Dump() {
	fmt.Printf("Number(%d)", n)
}

func (n Number) ToBoolean() Boolean {
	return n != 0
}

func (n Number) ToNumber() Number {
	return n
}

func (n Number) ToText() Text {
	return Text(strconv.Itoa(int(n)))
}

func (n Number) ToList() List {
	if n < 0 {
		panic("negative value given to list")
	}

	if n == 0 {
		return List{n}
	}

	// TODO: maybe this could be optimized?
	var list List
	for n != 0 {
		list = append(List{n % 10}, list...)
		n /= 10
	}

	return list
}
