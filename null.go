package knight

import "fmt"

type Null struct{}

func (n Null) Run() (Value, error) {
	return n, nil
}

func (n Null) Dump() {
	fmt.Print("Null()")
}

func (n Null) Bool() bool {
	return false
}

func (n Null) Int() int {
	return 0
}

func (n Null) String() string {
	return "null"
}
