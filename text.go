package knight

import "fmt"

type Text string

func (t Text) Run() (Value, error) {
	return t, nil
}

func (t Text) Dump() {
	fmt.Printf("Text(%s)", t)
}

func (t Text) Bool() bool {
	return t != ""
}

func (t Text) Int() int {
	var ret int

	fmt.Sscan(string(t), "%l", ret)

	return ret
}

func (t Text) String() string {
	return string(t)
}
