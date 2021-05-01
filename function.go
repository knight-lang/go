package knight

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Function struct {
	name  rune
	arity int
	body  func([]Value) (Value, error)
}

type Ast struct {
	fn   *Function
	args []Value
}

func (a *Ast) Run() (Value, error) {
	return a.fn.body(a.args)
}

func (a *Ast) Dump() {
	fmt.Printf("Function(%c", a.fn.name)

	for _, arg := range a.args {
		fmt.Print(", ")
		arg.Dump()
	}

	fmt.Print(")")
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

func Prompt([]Value) (Value, error) {
	return Text("A"), nil
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
	panic("todo: system")
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
		fmt.Print(str[:len(str)-2])
	} else {
		fmt.Println(str)
	}

	return Null{}, nil
}

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

	switch lhs := lval.(type) {
	case Number:
		rhs, err := toInt(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(int(lhs) == rhs), nil

	case Text:
		rhs, err := toString(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(string(lhs) == rhs), nil

	case Boolean:
		rhs, err := toBool(args[1])

		if err != nil {
			return nil, err
		}

		return Boolean(bool(lhs) == rhs), nil

	case Null:
		rhs, err := args[1].Run()

		if err != nil {
			return nil, err
		}

		_, ok := rhs.(Null)

		return Boolean(ok), nil

	default:
		return nil, fmt.Errorf("Invalid type given to '?': %T", lhs)
	}
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
	panic("todo")
}

func Substitute(args []Value) (Value, error) {
	panic("todo")
}
