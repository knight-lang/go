package knight

import "testing"


var args1 = []Value{Number(1), Number(2)}
var args2 = []Value{Text("1"), Text("2")}
var args3 = []Value{List{}, List{}}
func Benchmark1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := add1(args1); err != nil { panic("bad error [1]") }
		if _, err := add1(args2); err != nil { panic("bad error [2]") }
		if _, err := add1(args3); err != nil { panic("bad error [3]") }
	}
}

func Benchmark2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Add(args1); err != nil { panic("bad error") }
		if _, err := Add(args2); err != nil { panic("bad error") }
		if _, err := Add(args3); err != nil { panic("bad error") }
	}
}
