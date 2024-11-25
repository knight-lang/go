package knight

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
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

var (
	// this is a global variable because there's no way to read lines without a scanner.
	stdinScanner = bufio.NewScanner(os.Stdin)

	KnownFunctions = map[rune]*Function{
		'P': &Function{ name: 'P', arity: 0, fn: prompt },
		'R': &Function{ name: 'R', arity: 0, fn: random },
		'B': &Function{ name: 'B', arity: 1, fn: block },
		'C': &Function{ name: 'C', arity: 1, fn: call },
		'Q': &Function{ name: 'Q', arity: 1, fn: quit },
		'!': &Function{ name: '!', arity: 1, fn: not },
		'L': &Function{ name: 'L', arity: 1, fn: length },
		'D': &Function{ name: 'D', arity: 1, fn: dump },
		'O': &Function{ name: 'O', arity: 1, fn: output },
		'A': &Function{ name: 'A', arity: 1, fn: ascii },
		'~': &Function{ name: '~', arity: 1, fn: negate },
		',': &Function{ name: ',', arity: 1, fn: box },
		'[': &Function{ name: '[', arity: 1, fn: head },
		']': &Function{ name: ']', arity: 1, fn: tail },
		'+': &Function{ name: '+', arity: 2, fn: add },
		'-': &Function{ name: '-', arity: 2, fn: subtract },
		'*': &Function{ name: '*', arity: 2, fn: multiply },
		'/': &Function{ name: '/', arity: 2, fn: divide },
		'%': &Function{ name: '%', arity: 2, fn: remainder },
		'^': &Function{ name: '^', arity: 2, fn: exponentiate },
		'<': &Function{ name: '<', arity: 2, fn: lessthan },
		'>': &Function{ name: '>', arity: 2, fn: greaterthan },
		'?': &Function{ name: '?', arity: 2, fn: equalto },
		'&': &Function{ name: '&', arity: 2, fn: and },
		'|': &Function{ name: '|', arity: 2, fn: or },
		';': &Function{ name: ';', arity: 2, fn: then },
		'=': &Function{ name: '=', arity: 2, fn: assign },
		'W': &Function{ name: 'W', arity: 2, fn: while },
		'I': &Function{ name: 'I', arity: 3, fn: if_ },
		'G': &Function{ name: 'G', arity: 3, fn: get },
		'S': &Function{ name: 'S', arity: 4, fn: set },
	}
)

// initialize the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())	
}

func invalidType(fnName rune, value any) error {
	return fmt.Errorf("invalid type given to '%c': %T", fnName, value)
}

// runToString runs the `value` and then converts it to a string.
func runToString(value Value) (String, error) {
	ran, err := value.Run()
	if err != nil {
		return "", err
	}

	return ran.(Convertible).ToString(), nil
}

// runToInteger runs the `value` and then converts it to an integer.
func runToInteger(value Value) (Integer, error) {
	ran, err := value.Run()
	if err != nil {
		return 0, err
	}

	return ran.(Convertible).ToInteger(), nil
}

// runToInteger runs the `value` and then converts it to a boolean.
func runToBoolean(value Value) (Boolean, error) {
	ran, err := value.Run()
	if err != nil {
		return false, err
	}

	return ran.(Convertible).ToBoolean(), nil
}

// runToInteger runs the `value` and then converts it to a list.
func runToList(value Value) (List, error) {
	ran, err := value.Run()
	if err != nil {
		return nil, err
	}

	return ran.(Convertible).ToList(), nil
}

/** ARITY ZERO **/

// prompt reads a line from stdin, returning `Null` if we're closed.
func prompt(_ []Value) (Value, error) {
	if stdinScanner.Scan() {
		return String(strings.TrimRight(stdinScanner.Text(), "\r")), nil
	}

	if err := stdinScanner.Err(); err != nil {
		return nil, err
	}

	return Null{}, nil
}

// random returns a random `Integer`.
func random(_ []Value) (Value, error) {
	return Integer(rand.Int63()), nil
}

/** ARITY ONE **/

// box creates a list of its sole argument
func box(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return List{ran}, nil
}

// head returns the first element/character of a list/string.
//
// This returns an error if the argument is not a list or string, or is empty.
func head(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch container := ran.(type) {
	case List:
		if len(container) == 0 {
			return nil, errors.New("head on empty list")
		}
		return container[0], nil

	case String:
		if len(container) == 0 {
			return nil, errors.New("head on empty string")
		}
		return String(container[0]), nil

	default:
		return nil, invalidType('[', container)
	}
}

// tail returns a list/string of everything but the first element/char.
//
// This returns an error if the argument is not a list or string, or is empty.
func tail(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch container := ran.(type) {
	case List:
		if len(container) == 0 {
			return nil, errors.New("tail on empty list")
		}
		return container[1:], nil

	case String:
		if len(container) == 0 {
			return nil, errors.New("tail on empty string")
		}
		return container[1:], nil

	default:
		return nil, invalidType(']', container)
	}
}

// block returns its argument unevaluated.
func block(args []Value) (Value, error) {
	return args[0], nil
}

// call runs its argument twice.
func call(args []Value) (Value, error) {
	block, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return block.Run()
}

// quit exits the program with the given exit code.
func quit(args []Value) (Value, error) {
	exitStatus, err := runToInteger(args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(int(exitStatus))
	panic("<unreachable>")
}

// not returns the logical negation of its argument
func not(args []Value) (Value, error) {
	boolean, err := runToBoolean(args[0])
	if err != nil {
		return nil, err
	}

	return Boolean(!boolean), nil
}

// length converts its argument to a list, then returns its length.
//
// This returns an error if the argument is not a list or string.
func length(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	// There need to be two conditions, even though they're identical, as the `len` function is
	// operating on a different type (and so `case List, String:` wouldn't work).
	switch container := ran.(type) {
	case List:
		return Integer(len(container)), nil

	case String:
		return Integer(len(container)), nil

	default:
		return nil, invalidType('L', container)
	}
}

// dump prints a debugging representation of its argument to stdout, then returns it.
func dump(args []Value) (Value, error) {
	value, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	value.Dump()
	return value, nil
}

// output writes its argument to stdout, and returns null.
//
// If a `\` is the very last character, it's stripped and no newline is added. Otherwise, a newline
// is also printed.
func output(args []Value) (Value, error) {
	outputString, err := runToString(args[0])
	if err != nil {
		return nil, err
	}

	if outputString != "" && outputString[len(outputString) - 1] == '\\' {
		fmt.Print(outputString[:len(outputString) - 1])
	} else {
		fmt.Println(outputString)
	}

	return Null{}, nil
}

// ascii is essentially equivalent to `chr`/`ord` in other langauges, depending on its argument.
func ascii(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch ran := value.(type) {
	case Integer:
		if !utf8.ValidRune(rune(value)) {
			return nil, fmt.Errorf("invalid integer given to 'A': %d", value)
		}

		return String(rune(value)), nil

	case String:
		if value == "" {
			return nil, fmt.Errorf("empty string given to 'A'")
		}

		// We know the rune's not empty, so the panic's not a owrry.
		rune, _, _ := value.SplitFirstRune()
		return Integer(rune[0]), nil

	default:
		return nil, invalidType('A', value)
	}
}

// negate returns the numerical negation of its argument.
func negate(args []Value) (Value, error) {
	number, err := runToInteger(args[0])
	if err != nil {
		return nil, err
	}

	return -number, nil
}

/** ARITY TWO **/

// add adds two numbers/strings/lists together; it coerces the second argument.
func add(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}

		return lhs + rhs, nil

	case String:
		rhs, err := runToString(args[1])
		if err != nil {
			return nil, err
		}

		// using `strings.Builder` is a bit more efficient than concating and stuff.
		var sb strings.Builder
		sb.WriteString(string(lhs))
		sb.WriteString(string(rhs))

		return String(sb.String()), nil

	case List:
		rhs, err := runToList(args[1])
		if err != nil {
			return nil, err
		}

		return append(lhs, rhs...), nil

	default:
		return nil, fmt.Errorf("invalid type given to '+': %T", lhs)
	}
}

// subtract subtracts one number from another.
func subtract(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}

		return lhs - rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '-': %T", lval)
	}
}

// multiply multiplies two numbers, or repeats lists/strings; last argument's converted to a number.
func multiply(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}

		return lhs * rhs, nil

	case String:
		amount, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}
		if amount < 0 {
			return nil, fmt.Errorf("negative replication amount: %d", amount)
		}

		return String(strings.Repeat(string(lhs), int(amount))), nil

	case List:
		amount, err := runToInteger(args[1])
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

// divide divides the first argument by the second; errors out of the second is zero.
func divide(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
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

// remainder returns the remainder of `<arg1>/<arg2>`; errors out if second arg is zero.
func remainder(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}
		if rhs == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}

		return lhs % rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '%%': %T", lhs)
	}
}

// exponentiate raises the first argument to the power of the second, or joins lists. errors out on
// negative powers for integers.
func exponentiate(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lval.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}
		if rhs < 0 {
			return nil, fmt.Errorf("Exponentiation of negative power attempted")
		}

		// All 32 bit number exponentiations that can be represented in 32 bits can be done with
		// 64 bit floats and a "powf" function.
		return Integer(math.Pow(float64(lhs), float64(rhs))), nil

	case List:
		sep, err := runToString(args[1])
		if err != nil {
			return nil, err
		}

		return String(lhs.Join(string(sep))), nil

	default:
		return nil, fmt.Errorf("invalid type given to '^': %T", lhs)
	}
}

func compare(lhs, rhs Value, fn rune) (int, error) {
	switch lhs := lhs.(type) {
	case Integer:
		return int(lhs - rhs.(Convertible).ToInteger()), nil

	case String:
		return strings.Compare(string(lhs), string(rhs.(Convertible).ToString())), nil

	case Boolean:
		return int(lhs.ToInteger() - rhs.(Convertible).ToBoolean().ToInteger()), nil

	case List:
		rhs := rhs.(Convertible).ToList()
		minLen := len(lhs)
		if len(rhs) < minLen {
			minLen = len(rhs)
		}

		for i := 0; i < minLen; i++ {
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

// lessThan returns whether the first argument is less than the second.
func lessthan(args []Value) (Value, error) {
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

// greaterThan returns whether the first argument is greater than the second.
func greaterthan(args []Value) (Value, error) {
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

// equalTo returns whether its two arguments are equal to one other.
func equalto(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rval, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	// `DeepEqual` happens to correspond exactly to Knight's equality semantics
	return Boolean(reflect.DeepEqual(lval, rval)), nil
}

// and evaluates the first argument, then either returns that if it's truthy, or otherwise evaluates
// and returns the second argument.
func and(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	if lval.(Convertible).ToBoolean() {
		return args[1].Run()
	}

	return lval, nil
}

// or evaluates the first argument, then either returns that if it's falsey, or otherwise evaluates
// and returns the second argument.
func or(args []Value) (Value, error) {
	lval, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	if !lval.(Convertible).ToBoolean() {
		return args[1].Run()
	}

	return lval, nil
}

// then evaluates the first argument, then evaluates and returns the second argument.
func then(args []Value) (Value, error) {
	if _, err := args[0].Run(); err != nil {
		return nil, err
	}

	return args[1].Run()
}

// assign assigns the second argument to the first argument (which must be a `Variable`).
func assign(args []Value) (Value, error) {
	variable, ok := args[0].(*Variable)
	if !ok {
		return nil, fmt.Errorf("invalid type given to '=': %T", args[0])
	}

	value, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	variable.Assign(value)

	return value, nil
}

// while evaluates the second argument whilst the first is true.
func while(args []Value) (Value, error) {
	for {
		cond, err := runToBoolean(args[0])
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

/** ARITY THREE **/

// if will evaluate and return either the 2nd or 3rd argument, depending on the 1st's truthiness
func if_(args []Value) (Value, error) {
	cond, err := runToBoolean(args[0])
	if err != nil {
		return nil, err
	}

	if cond {
		return args[1].Run()
	}

	return args[2].Run()
}

// get returns a sublist/string with start and length of the second and third elements.
func get(args []Value) (Value, error) {
	collection, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	start, err := runToInteger(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to GET (%d)", start)
	}

	length, err := runToInteger(args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to GET (%d)", length)
	}

	switch collection := collection.(type) {
	case String:
		if Integer(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		return collection[start : start+length], nil

	case List:
		if Integer(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		return collection[start : start+length], nil

	default:
		return nil, fmt.Errorf("invalid type given to 'G': %T", collection)
	}
}

/** ARITY FOUR **/

// set returns a list/string where the range `[<arg2>, <arg2>+<arg3>)` is replaced by the fourth.
func set(args []Value) (Value, error) {
	collection, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	start, err := runToInteger(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to SET (%d)", start)
	}

	length, err := runToInteger(args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to SET (%d)", length)
	}

	switch collection := collection.(type) {
	case String:
		if Integer(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		replacement, err := runToString(args[3])
		if err != nil {
			return nil, err
		}

		return collection[:start] + replacement + collection[start+length:], nil

	case List:
		if Integer(len(collection)) < start+length {
			return nil, fmt.Errorf("len (%d) < start (%d) + len (%d)", len(collection), start, length)
		}

		begin := collection[:start]
		end := collection[start+length:]

		middle, err := runToList(args[3])
		if err != nil {
			return nil, err
		}

		ret := make(List, 0, len(collection)-int(length)+len(middle))
		return append(append(append(ret, begin...), middle...), end...), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'S': %T", collection)

	}
}
