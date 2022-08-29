package knight

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Function struct {
	name  rune
	arity int
	body  func([]Value) (Value, error)
}

type Ast struct {
	fun  *Function
	args []Value
}

var functions map[rune]*Function = make(map[rune]*Function)

func GetFunction(r rune) *Function {
	val, ok := functions[r]

	if !ok {
		return nil
	}

	return val
}

func RegisterFunction(name rune, arity int, body func([]Value) (Value, error)) {
	functions[name] = &Function{name: name, arity: arity, body: body}
}

func (a *Ast) Run() (Value, error) {
	return a.fun.body(a.args)
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

func toString(value Value) (string, error) {
	ran, err := value.Run()

	if err != nil {
		return "", err
	}

	return ran.(Literal).String(), nil
}

func toInt(value Value) (int, error) {
	ran, err := value.Run()

	if err != nil {
		return 0, err
	}

	return ran.(Literal).Int(), nil
}

func toBool(value Value) (bool, error) {
	ran, err := value.Run()

	if err != nil {
		return false, err
	}

	return ran.(Literal).Bool(), nil
}

func toList(value Value) ([]Value, error) {
	ran, err := value.Run()

	if err != nil {
		return nil, err
	}

	return ran.(Literal).List(), nil
}

/** ARITY ZERO **/

var reader = bufio.NewReader(os.Stdin)

func Prompt([]Value) (Value, error) {
	line, _, err := reader.ReadLine()

	if err != nil {
		return nil, err
	}

	return Text(line), err
}

func Random([]Value) (Value, error) {
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
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return ran.Run()
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

	command := exec.Command(shell, "-c", cmd)
	command.Stdin = os.Stdin
	stdout, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("unable to read command result: %s", err)
	}

	return Text(stdout), nil
}

func Quit(args []Value) (Value, error) {
	code, err := toInt(args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(code)
	panic("unreachable")
}

func Not(args []Value) (Value, error) {
	bool, err := toBool(args[0])
	if err != nil {
		return nil, err
	}

	return Boolean(!bool), nil
}

func Length(args []Value) (Value, error) {
	list, err := toList(args[0])
	if err != nil {
		return nil, err
	}

	return Number(len(list)), nil
}

func Dump(args []Value) (Value, error) {
	val, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	val.Dump()
	fmt.Println()

	return val, nil
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
	val, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := val.(type) {
	case Number:
		if lhs <= 0 || 127 < lhs {
			return nil, fmt.Errorf("invalid number given to 'A': %d", lhs)
		}

		return Text(rune(lhs)), nil

	case Text:
		if lhs == "" {
			return nil, fmt.Errorf("empty string given to 'A'")
		}

		return Number(lhs[0]), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'A': %T", lhs)
	}
}

/** ARITY TWO **/

func Add(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toInt(args[1])
		if err != nil {
			return nil, err
		}

		return lhs + Number(rhs), nil

	case Text:
		rhs, err := toString(args[1])
		if err != nil {
			return nil, err
		}

		var sb strings.Builder
		sb.WriteString(string(lhs))
		sb.WriteString(rhs)

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
		rhs, err := toInt(args[1])
		if err != nil {
			return nil, err
		}

		return lhs - Number(rhs), nil

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
		rhs, err := toInt(args[1])
		if err != nil {
			return nil, err
		}

		return lhs * Number(rhs), nil

	case Text:
		amount, err := toInt(args[1])
		if err != nil {
			return nil, err
		}
		if amount < 0 {
			return nil, fmt.Errorf("negative replication amount: %d", amount)
		}

		return Text(strings.Repeat(string(lhs), amount)), nil

	case List:
		amount, err := toInt(args[1])
		if err != nil {
			return nil, err
		}
		if amount < 0 {
			return nil, fmt.Errorf("negative replication amount: %d", amount)
		}

		slice := make(List, 0, len(lhs)*amount)

		for i := 0; i < amount; i++ {
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
		rhs, err := toInt(args[1])
		if err != nil {
			return nil, err
		}
		if rhs == 0 {
			return nil, fmt.Errorf("division by zero")
		}

		return lhs / Number(rhs), nil

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
		rhs, err := toInt(args[1])
		if err != nil {
			return nil, err
		}
		if rhs == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}

		return lhs % Number(rhs), nil

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
		rhs, err := toInt(args[1])
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

		var sb strings.Builder

		for i := 0; i < len(lhs); i++ {
			if i != 0 {
				sb.WriteString(sep)
			}

			sb.WriteString(lhs[i].(Literal).String())
		}

		return Text(sb.String()), nil

	default:
		return nil, fmt.Errorf("invalid type given to '^': %T", lhs)
	}
}

func compare(lhs, rhs Value, fn rune) (int, error) {
	switch lhs := lhs.(type) {
	case Number:
		return int(lhs) - rhs.(Literal).Int(), nil

	case Text:
		return strings.Compare(string(lhs), rhs.(Literal).String()), nil

	case Boolean:
		conv := func(b bool) int {
			if b {
				return 1
			}
			return 0
		}

		return conv(bool(lhs)) - conv(rhs.(Literal).Bool()), nil
	case List:
		r := rhs.(Literal).List()
		for i := 0; i < int(math.Min(float64(len(lhs)), float64(len(r)))); i++ {
			cmp, err := compare(lhs[i], r[i], fn)
			if err != nil {
				return 0, err
			}
			if cmp != 0 {
				return cmp, nil
			}
		}
		return len(lhs) - len(r), nil

	default:
		return 0, fmt.Errorf("invalid type given to %q: %T", fn, lhs)
	}
}

func LessThan(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rval, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	cmp, err := compare(lval, rval, '<')
	if err != nil {
		return nil, err
	}

	return Boolean(cmp < 0), nil
}

func GreaterThan(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rval, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	cmp, err := compare(lval, rval, '>')
	if err != nil {
		return nil, err
	}

	return Boolean(cmp > 0), nil
}

func equals(lhs, rhs Value) bool {
	if lhs == rhs {
		return true
	}

	switch lhs := lhs.(type) {
	case Number:
		if rhs, ok := rhs.(Number); ok {
			return lhs == rhs
		}
		return false

	case Text:
		if rhs, ok := rhs.(Text); ok {
			return lhs == rhs
		}
		return false

	case Boolean:
		if rhs, ok := rhs.(Boolean); ok {
			return lhs == rhs
		}
		return false

	case Null:
		_, ok := rhs.(Null)
		return ok

	case List:
		rhs, ok := rhs.(List)
		if !ok {
			return false
		}
		if len(lhs) != len(rhs) {
			return false
		}

		for i := 0; i < len(lhs); i++ {
			if !equals(lhs[i], rhs[i]) {
				return false
			}
		}

		return true

	default:
		return false
	}
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

	return Boolean(equals(lval, rval)), nil
}

func And(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	if lval.(Literal).Bool() {
		return args[1].Run()
	}

	return lval, nil
}

func Or(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	if !lval.(Literal).Bool() {
		return args[1].Run()
	}

	return lval, nil
}

func Then(args []Value) (Value, error) {
	_, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return args[1].Run()
}

func Assign(args []Value) (Value, error) {
	lval, ok := args[0].(*Variable)
	if !ok {
		return nil, fmt.Errorf("invalid type given to '=': %T", lval)
	}

	rval, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	lval.Assign(rval)

	return rval, nil
}

func While(args []Value) (Value, error) {
	for {
		cond, err := toBool(args[0])
		if err != nil {
			return nil, err
		}

		if !cond {
			return Null{}, nil
		}

		if _, err = args[1].Run(); err != nil {
			return nil, err
		}
	}

}

func Range(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Number:
		start := lhs
		stop, err := toInt(args[1])
		if err != nil {
			return nil, err
		}

		if start > Number(stop) {
			return nil, fmt.Errorf("invalid values to range: %d > %d", start, stop)
		}

		ret := make(List, 0, Number(stop)-start)
		for curr := start; curr != Number(stop); curr++ {
			ret = append(ret, curr)
		}
		return ret, nil

	case Text:
		if len(lhs) == 0 {
			return nil, fmt.Errorf("empty start given to range")
		}
		start := lhs[0]

		rhs, err := toString(args[1])
		if err != nil {
			return nil, err
		}
		if len(rhs) == 0 {
			return nil, fmt.Errorf("empty stop given to range")
		}
		stop := rhs[0]

		ret := make(List, 0, int32(stop)-int32(start))
		for curr := start; curr != stop; curr++ {
			ret = append(ret, Text(curr))
		}
		return ret, nil

	default:
		return nil, fmt.Errorf("invalid type given to '.': %T", lhs)
	}
}

/** ARITY THREE **/

func If(args []Value) (Value, error) {
	cond, err := toBool(args[0])
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

	start, err := toInt(args[1])
	if err != nil {
		return nil, err
	}

	amnt, err := toInt(args[2])
	if err != nil {
		return nil, err
	}

	switch lhs := collection.(type) {
	case Text:
		if len(lhs) <= start+amnt {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(lhs), start, amnt)
		}

		return lhs[start : start+amnt], nil

	default:
		return nil, fmt.Errorf("invalid type given to '.': %T", lhs)
	}
}

/** ARITY FOUR **/

func Substitute(args []Value) (Value, error) {
	str, err := toString(args[0])

	if err != nil {
		return nil, err
	}

	start, err := toInt(args[1])

	if err != nil {
		return nil, err
	}

	amnt, err := toInt(args[2])

	if err != nil {
		return nil, err
	}

	repl, err := toString(args[3])

	if err != nil {
		return nil, err
	}

	if start == len(str) && amnt == 0 {
		return Text(str + repl), nil
	}

	if start == 0 && len(repl) == 0 {
		return Text(str[amnt:]), nil
	}

	return Text(str[:start] + repl + str[start+amnt:]), nil
}
