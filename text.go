package knight

import "fmt"

type Text string

func (t Text) Run() (Value, error) {
	return t, nil
}

func (t Text) Dump() {
	fmt.Printf("String(%s)", t)
}

func (t Text) Bool() bool {
	return t != ""
}

func (t Text) Int() int {
	var ret int
	fmt.Sscanf(string(t), "%d", &ret)
	return ret
}

func (t Text) String() string {
	return string(t)
}

func (t Text) List() []Value {
	list := make([]Value, 0, len(t))

	for i := 0; i < len(t); i++ {
		list[i] = Text(string(t[i]))
	}

	return list
}
