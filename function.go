package knight

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"bufio"
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

var functions map[rune] *Function = make(map[rune] *Function)

func GetFunction(r rune) *Function {
	val, ok := functions[r]

	if !ok {
		return nil
	}

	return val
}

func RegisterFunction(name rune, arity int, body func([]Value) (Value, error)) {
	functions[name] = &Function { name: name, arity: arity, body: body }
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

	RegisterFunction('E', 1, Eval)
	RegisterFunction('B', 1, Block)
	RegisterFunction('C', 1, Call)
	RegisterFunction('`', 1, System)
	RegisterFunction('Q', 1, Quit)
	RegisterFunction('!', 1, Not)
	RegisterFunction('L', 1, Length)
	RegisterFunction('D', 1, Dump)
	RegisterFunction('O', 1, Output)

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

func Eval(args []Value) (Value, error) {
	str, err := toString(args[0])

	if err != nil {
		return nil, err
	}

	return Run(str)
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

	return nil, nil
}

func Not(args []Value) (Value, error) {
	bool, err := toBool(args[0])

	if err != nil {
		return nil, err
	}

	return Boolean(!bool), nil
}

func Length(args []Value) (Value, error) {
	str, err := toString(args[0])

	if err != nil {
		return nil, err
	}

	return Number(len(str)), nil
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

		return Number(int(lhs) + rhs), nil

	case Text:
		rhs, err := toString(args[1])

		if err != nil {
			return nil, err
		}

		return Text(string(lhs) + rhs), nil

	default:
		return nil, fmt.Errorf("Invalid type given to '+': %T", lhs)
	}
}

func Subtract(args []Value) (Value, error) {
	lval, err := args[0].Run()

	if err != nil {
		return nil, err
	}

	lhs, ok := lval.(Number)

	if !ok {
		return nil, fmt.Errorf("Invalid type given to '-': %T", lval)
	}

	rhs, err := toInt(args[1])

	if err != nil {
		return nil, err
	}

	return Number(int(lhs) - rhs), nil
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

		return Number(int(lhs) * rhs), nil

	case Text:
		rhs, err := toInt(args[1])

		if err != nil {
			return nil, err
		}

		if rhs < 0 {
			return nil, fmt.Errorf("Negative replication amount: %d", rhs)
		}

		return Text(strings.Repeat(string(lhs), rhs)), nil

	default:
		return nil, fmt.Errorf("Invalid type given to '*': %T", lhs)
	}
}

func Divide(args []Value) (Value, error) {
	lval, err := args[0].Run()

	if err != nil {
		return nil, err
	}

	lhs, ok := lval.(Number)

	if !ok {
		return nil, fmt.Errorf("Invalid type given to '/': %T", lval)
	}

	rhs, err := toInt(args[1])

	if err != nil {
		return nil, err
	}

	if rhs == 0 {
		return nil, fmt.Errorf("Division by zero attempted")
	}

	return Number(int(lhs) / rhs), nil
}

func Modulo(args []Value) (Value, error) {
	lval, err := args[0].Run()

	if err != nil {
		return nil, err
	}

	lhs, ok := lval.(Number)

	if !ok {
		return nil, fmt.Errorf("Invalid type given to '%': %T", lval)
	}

	rhs, err := toInt(args[1])

	if err != nil {
		return nil, err
	}

	if rhs == 0 {
		return nil, fmt.Errorf("Modulo by zero attempted")
	}

	return Number(int(lhs) % rhs), nil
}

func Exponentiate(args []Value) (Value, error) {
	lval, err := args[0].Run()

	if err != nil {
		return nil, err
	}

	lhs, ok := lval.(Number)

	if !ok {
		return nil, fmt.Errorf("Invalid type given to '^': %T", lval)
	}

	rhs, err := toInt(args[1])

	if err != nil {
		return nil, err
	}

	if lhs == 0 && rhs < 0 {
		return nil, fmt.Errorf("Exponentiation of zero to a negative power attempted")
	}

	return Number(int(math.Pow(float64(int(lhs)), float64(rhs)))), nil

}

func LessThan(args []Value) (Value, error) {
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

		return Boolean(int(lhs) < rhs), nil

	case Text:
		rhs, err := toString(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(string(lhs) < rhs), nil

	case Boolean:
		rhs, err := toBool(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(!bool(lhs) && rhs), nil

	default:
		return nil, fmt.Errorf("Invalid type given to '<': %T", lhs)
	}
}

func GreaterThan(args []Value) (Value, error) {
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

		return Boolean(int(lhs) > rhs), nil

	case Text:
		rhs, err := toString(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(string(lhs) > rhs), nil

	case Boolean:
		rhs, err := toBool(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(bool(lhs) && !rhs), nil

	default:
		return nil, fmt.Errorf("Invalid type given to '>': %T", lhs)
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

	return Boolean(lval == rval), nil
}

func And(args []Value) (Value, error) {
	lval, err := args[0].Run()

	if err != nil {
		return nil, err
	}

	lhs, ok := lval.(Literal)

	if !ok {
		return nil, fmt.Errorf("Invalid type given to '&': %T", lval)
	}

	if lhs.Bool() {
		return args[1].Run()
	}

	return lval, nil
}

func Or(args []Value) (Value, error) {
	lval, err := args[0].Run()

	if err != nil {
		return nil, err
	}

	lhs, ok := lval.(Literal)

	if !ok {
		return nil, fmt.Errorf("Invalid type given to '|': %T", lval)
	}

	if !lhs.Bool() {
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
		return nil, fmt.Errorf("Invalid type given to '=': %T", lval)
	}

	rval, err := args[1].Run();

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
			break
		}

		_, err = args[1].Run();

		if err != nil {
			return nil, err
		}
	}

	return Null{}, nil
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
	str, err := toString(args[0])

	if err != nil {
		return nil, err
	}

	start, err := toInt(args[1])

	if err != nil  {
		return nil, err
	}

	amnt, err := toInt(args[2])

	if err != nil {
		return nil, err
	}

	if start == len(str) {
		return Text(""), nil
	}


	return Text(str[start:start + amnt]), nil
}

/** ARITY FOUR **/

func Substitute(args []Value) (Value, error) {
	str, err := toString(args[0])

	if err != nil {
		return nil, err
	}

	start, err := toInt(args[1])

	if err != nil  {
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

	return Text(str[:start] + repl + str[start+amnt:]), nil
}
