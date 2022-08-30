package knight

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

// Function is a type representing a function within Knight.
//
// Each `Function`'s `fn` expects to receive exactly `arity` arguments
type Function struct {
	name  rune
	arity int
	fn    func([]Value) (Value, error)
}

var functions map[rune]*Function = make(map[rune]*Function)

func NewFunction(name rune, arity int, fn func([]Value) (Value, error)) *Function {
	return &Function{name: name, arity: arity, fn: fn}
}

func GetFunction(name rune) *Function {
	if fn, ok := functions[name]; ok {
		return fn
	}

	return nil
}

func RegisterFunction(name rune, arity int, fn func([]Value) (Value, error)) {
	functions[name] = &Function{name: name, arity: arity, fn: fn}
}

type Ast struct {
	fun  *Function
	args []Value
}

func (a *Ast) Run() (Value, error) {
	return a.fun.fn(a.args)
}

func (a *Ast) Dump() {
	fmt.Printf("Function(%c", a.fun.name)

	for _, arg := range a.args {
		fmt.Print(", ")
		arg.Dump()
	}

	fmt.Print(")")
}

func init() {
	rand.Seed(time.Now().UnixNano())

	RegisterFunction('P', 0, Prompt)
	RegisterFunction('R', 0, Random)

	RegisterFunction(',', 1, Box)
	RegisterFunction('B', 1, Block)
	RegisterFunction('C', 1, Call)
	RegisterFunction('`', 1, System)
	RegisterFunction('Q', 1, Quit)
	RegisterFunction('!', 1, Not)
	RegisterFunction('L', 1, Length)
	RegisterFunction('D', 1, Dump)
	RegisterFunction('O', 1, Output)
	RegisterFunction('A', 1, Ascii)
	RegisterFunction('~', 1, Negate)

	RegisterFunction('+', 2, Add)
	RegisterFunction('-', 2, Subtract)
	RegisterFunction('*', 2, Multiply)
	RegisterFunction('/', 2, Divide)
	RegisterFunction('%', 2, Modulo)
	RegisterFunction('^', 2, Exponentiate)
	RegisterFunction('<', 2, LessThan)
	RegisterFunction('>', 2, GreaterThan)
	RegisterFunction('?', 2, EqualTo)
	RegisterFunction('&', 2, And)
	RegisterFunction('|', 2, Or)
	RegisterFunction(';', 2, Then)
	RegisterFunction('=', 2, Assign)
	RegisterFunction('W', 2, While)
	RegisterFunction('.', 2, Range)

	RegisterFunction('I', 3, If)
	RegisterFunction('G', 3, Get)

	RegisterFunction('S', 4, Substitute)
}

func toString(value Value) (Text, error) {
	ran, err := value.Run()

	if err != nil {
		return "", err
	}

	return ran.(Literal).ToText(), nil
}

func toNumber(value Value) (Number, error) {
	ran, err := value.Run()

	if err != nil {
		return Number(0), err
	}

	return ran.(Literal).ToNumber(), nil
}

func toBoolean(value Value) (Boolean, error) {
	ran, err := value.Run()

	if err != nil {
		return false, err
	}

	return ran.(Literal).ToBoolean(), nil
}

func toList(value Value) (List, error) {
	ran, err := value.Run()

	if err != nil {
		return nil, err
	}

	return ran.(Literal).ToList(), nil
}

/** ARITY ZERO **/

var reader = bufio.NewReader(os.Stdin)

func Prompt(_ []Value) (Value, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	return Text(line), err
}

func Random(_ []Value) (Value, error) {
	return Number(rand.Int63()), nil
}

/** ARITY ONE **/
func Box(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return List{ran}, nil
}

func Block(args []Value) (Value, error) {
	return args[0], nil
}

func Call(args []Value) (Value, error) {
	block, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return block.Run()
}

func System(args []Value) (Value, error) {
	cmd, err := toString(args[0])
	if err != nil {
		return nil, err
	}

	shell := "/bin/sh"
	if s := os.Getenv("SHELL"); s != "" {
		shell = s
	}

	command := exec.Command(shell, "-c", string(cmd))
	command.Stdin = os.Stdin
	stdout, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("unable to read command result: %s", err)
	}

	return Text(stdout), nil
}

func Quit(args []Value) (Value, error) {
	code, err := toNumber(args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(int(code))
	panic("unreachable")
}

func Not(args []Value) (Value, error) {
	boolean, err := toBoolean(args[0])
	if err != nil {
		return nil, err
	}

	return !boolean, nil
}

func Length(args []Value) (Value, error) {
	list, err := toList(args[0])
	if err != nil {
		return nil, err
	}

	return Number(len(list)), nil
}

func Dump(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	ran.Dump()
	fmt.Println()

	return ran, nil
}

func Output(args []Value) (Value, error) {
	str, err := toString(args[0])
	if err != nil {
		return nil, err
	}

	if str != "" && str[len(str)-1] == '\\' {
		fmt.Print(str[:len(str)-1])
	} else {
		fmt.Println(str)
	}

	return Null{}, nil
}

func Ascii(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch value := ran.(type) {
	case Number:
		if !utf8.ValidRune(rune(value)) {
			return nil, fmt.Errorf("invalid integer given to 'A': %d", value)
		}

		return Text(rune(value)), nil

	case Text:
		if value == "" {
			return nil, fmt.Errorf("empty string given to 'A'")
		}

		rune, _ := value.FirstRune()
		return Number(rune), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'A': %T", value)
	}
}

func Negate(args []Value) (Value, error) {
	number, err := toNumber(args[0])
	if err != nil {
		return nil, err
	}

	return -number, nil
}

/** ARITY TWO **/

func Add(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}

		return lhs + rhs, nil

	case Text:
		rhs, err := toString(args[1])
		if err != nil {
			return nil, err
		}

		// using `strings.Builder` is a bit more efficient than concating and stuff.
		var sb strings.Builder
		sb.WriteString(string(lhs))
		sb.WriteString(string(rhs))

		return Text(sb.String()), nil

	case List:
		rhs, err := toList(args[1])
		if err != nil {
			return nil, err
		}

		return append(lhs, rhs...), nil

	default:
		return nil, fmt.Errorf("invalid type given to '+': %T", lhs)
	}
}

func Subtract(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}

		return lhs - rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '-': %T", lval)
	}
}

func Multiply(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}

		return lhs * rhs, nil

	case Text:
		amount, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}
		if amount < 0 {
			return nil, fmt.Errorf("negative replication amount: %d", amount)
		}

		return Text(strings.Repeat(string(lhs), int(amount))), nil

	case List:
		amount, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}
		if amount < 0 {
			return nil, fmt.Errorf("negative replication amount: %d", amount)
		}

		slice := make(List, 0, len(lhs)*int(amount))

		for i := 0; i < int(amount); i++ {
			slice = append(slice, lhs...)
		}

		return slice, nil

	default:
		return nil, fmt.Errorf("invalid type given to '*': %T", lhs)
	}
}

func Divide(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}
		if rhs == 0 {
			return nil, fmt.Errorf("division by zero")
		}

		return lhs / rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '/': %T", lhs)
	}
}

func Modulo(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}
		if rhs == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}

		return lhs % rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '%': %T", lhs)
	}

}

func Exponentiate(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}
		if rhs < 0 {
			return nil, fmt.Errorf("Exponentiation of negative power attempted")
		}

		return Number(math.Pow(float64(lhs), float64(rhs))), nil

	case List:
		sep, err := toString(args[1])
		if err != nil {
			return nil, err
		}

		return Text(lhs.Join(string(sep))), nil

	default:
		return nil, fmt.Errorf("invalid type given to '^': %T", lhs)
	}
}

func compare(lhs, rhs Value, fn rune) (int, error) {
	switch lhs := lhs.(type) {
	case Number:
		return int(lhs - rhs.(Literal).ToNumber()), nil

	case Text:
		return strings.Compare(string(lhs), string(rhs.(Literal).ToText())), nil

	case Boolean:
		return int(lhs.ToNumber() - rhs.(Literal).ToBoolean().ToNumber()), nil

	case List:
		rhs := rhs.(Literal).ToList()
		min_len := len(lhs)
		if len(rhs) < min_len {
			min_len = len(rhs)
		}

		for i := 0; i < min_len; i++ {
			cmp, err := compare(lhs[i], rhs[i], fn)
			if err != nil {
				return 0, err
			}

			if cmp != 0 {
				return cmp, nil
			}
		}

		return len(lhs) - len(rhs), nil

	default:
		return 0, fmt.Errorf("invalid type given to %q: %T", fn, lhs)
	}
}

func LessThan(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rhs, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	cmp, err := compare(lhs, rhs, '<')
	if err != nil {
		return nil, err
	}

	return Boolean(cmp < 0), nil
}

func GreaterThan(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rhs, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	cmp, err := compare(lhs, rhs, '>')
	if err != nil {
		return nil, err
	}

	return Boolean(cmp > 0), nil
}

func EqualTo(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rval, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	// `DeepEqual` happens to correspond exactly to knight's equality semantics
	return Boolean(reflect.DeepEqual(lval, rval)), nil
}

func And(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	if lval.(Literal).ToBoolean() {
		return args[1].Run()
	}

	return lval, nil
}

func Or(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	if !lval.(Literal).ToBoolean() {
		return args[1].Run()
	}

	return lval, nil
}

func Then(args []Value) (Value, error) {
	if _, err := args[0].Run(); err != nil {
		return nil, err
	}

	return args[1].Run()
}

func Assign(args []Value) (Value, error) {
	variable, ok := args[0].(*Variable)
	if !ok {
		return nil, fmt.Errorf("invalid type given to '=': %T", variable)
	}

	value, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	variable.Assign(value)

	return value, nil
}

func While(args []Value) (Value, error) {
	for {
		cond, err := toBoolean(args[0])
		if err != nil {
			return nil, err
		}

		if !cond {
			break
		}

		if _, err = args[1].Run(); err != nil {
			return nil, err
		}
	}

	return Null{}, nil
}

func Range(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		start := lhs
		stop, err := toNumber(args[1])
		if err != nil {
			return nil, err
		}

		if stop < start {
			return nil, fmt.Errorf("invalid values to range: %d > %d", start, stop)
		}

		list := make(List, 0, stop-start)
		for current := start; current < stop; current++ {
			list = append(list, current)
		}

		return list, nil

	case Text:
		if lhs == "" {
			return nil, fmt.Errorf("empty start given to range")
		}

		rhs, err := toString(args[1])
		if err != nil {
			return nil, err
		}
		if rhs == "" {
			return nil, fmt.Errorf("empty stop given to range")
		}

		start, _ := lhs.FirstRune()
		stop, _ := rhs.FirstRune()

		if stop < start {
			return nil, fmt.Errorf("invalid values to range: %q > %q", start, stop)
		}

		rng := make(List, 0, stop-start)
		for curr := start; curr != stop; curr++ {
			if utf8.ValidRune(curr) {
				rng = append(rng, Text(curr))
			}
		}

		return rng, nil

	default:
		return nil, fmt.Errorf("invalid type given to '.': %T", lhs)
	}
}

/** ARITY THREE **/

func If(args []Value) (Value, error) {
	cond, err := toBoolean(args[0])
	if err != nil {
		return nil, err
	}

	if cond {
		return args[1].Run()
	}

	return args[2].Run()
}

func Get(args []Value) (Value, error) {
	collection, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	start, err := toNumber(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to GET (%d)", start)
	}

	length, err := toNumber(args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to GET (%d)", length)
	}

	switch collection := collection.(type) {
	case Text:
		if Number(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		return collection[start : start+length], nil

	case List:
		if Number(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		// Special case for returning _just_ the element at that index.
		if length == 0 {
			return collection[start], nil
		}

		return collection[start : start+length], nil

	default:
		return nil, fmt.Errorf("invalid type given to 'G': %T", collection)
	}
}

/** ARITY FOUR **/

func Substitute(args []Value) (Value, error) {
	collection, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	start, err := toNumber(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to SET (%d)", start)
	}

	length, err := toNumber(args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to SET (%d)", length)
	}

	switch collection := collection.(type) {
	case Text:
		if Number(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		replacement, err := toString(args[3])
		if err != nil {
			return nil, err
		}

		return collection[:start] + replacement + collection[start+length:], nil

	case List:
		if Number(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		begin := collection[:start]
		end := collection[start+length:]

		middle, err := toList(args[3])
		if err != nil {
			return nil, err
		}

		ret := make(List, 0, len(collection)-int(length)+len(middle))
		return append(append(append(ret, begin...), middle...), end...), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'S': %T", collection)

	}
}
