package knight

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"slices"
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
		// Arity 0
		'T': &Function{name: 'T', arity: 0, fn: true_},
		'F': &Function{name: 'F', arity: 0, fn: false_},
		'N': &Function{name: 'N', arity: 0, fn: null},
		'@': &Function{name: '@', arity: 0, fn: emptyList},
		'P': &Function{name: 'P', arity: 0, fn: prompt},
		'R': &Function{name: 'R', arity: 0, fn: random},

		// Arity 1
		':': &Function{name: ':', arity: 1, fn: noop},
		'B': &Function{name: 'B', arity: 1, fn: block},
		'C': &Function{name: 'C', arity: 1, fn: call},
		'Q': &Function{name: 'Q', arity: 1, fn: quit},
		'!': &Function{name: '!', arity: 1, fn: not},
		'L': &Function{name: 'L', arity: 1, fn: length},
		'D': &Function{name: 'D', arity: 1, fn: dump},
		'O': &Function{name: 'O', arity: 1, fn: output},
		'A': &Function{name: 'A', arity: 1, fn: ascii},
		'~': &Function{name: '~', arity: 1, fn: negate},
		',': &Function{name: ',', arity: 1, fn: box},
		'[': &Function{name: '[', arity: 1, fn: head},
		']': &Function{name: ']', arity: 1, fn: tail},

		// Arity 2
		'+': &Function{name: '+', arity: 2, fn: add},
		'-': &Function{name: '-', arity: 2, fn: subtract},
		'*': &Function{name: '*', arity: 2, fn: multiply},
		'/': &Function{name: '/', arity: 2, fn: divide},
		'%': &Function{name: '%', arity: 2, fn: remainder},
		'^': &Function{name: '^', arity: 2, fn: exponentiate},
		'<': &Function{name: '<', arity: 2, fn: lessThan},
		'>': &Function{name: '>', arity: 2, fn: greaterThan},
		'?': &Function{name: '?', arity: 2, fn: equalTo},
		'&': &Function{name: '&', arity: 2, fn: and},
		'|': &Function{name: '|', arity: 2, fn: or},
		';': &Function{name: ';', arity: 2, fn: then},
		'=': &Function{name: '=', arity: 2, fn: assign},
		'W': &Function{name: 'W', arity: 2, fn: while},

		// Arity 3
		'I': &Function{name: 'I', arity: 3, fn: if_},
		'G': &Function{name: 'G', arity: 3, fn: get},

		// Arity 4
		'S': &Function{name: 'S', arity: 4, fn: set},
	}
)

// initialize the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

/** ARITY ZERO **/
func true_(_ []Value) (Value, error) {
	return Boolean(true), nil
}

func false_(_ []Value) (Value, error) {
	return Boolean(false), nil
}

func null(_ []Value) (Value, error) {
	return &Null{}, nil
}

func emptyList(_ []Value) (Value, error) {
	return &List{}, nil
}

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

// noop simply executes its only argument and returns it
func noop(args []Value) (Value, error) {
	return args[0].Run()
}

// box creates a list of its sole argument
func box(args []Value) (Value, error) {
	value, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return List{value}, nil
}

// head returns the first element/character of a list/string.
//
// This returns an error if the argument is not a list or string, or is empty.
func head(args []Value) (Value, error) {
	container, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch container := container.(type) {
	case List:
		if len(container) == 0 {
			return nil, errors.New("empty list given to '['")
		}

		return container[0], nil

	case String:
		if len(container) == 0 {
			return nil, errors.New("empty string given to '['")
		}

		return String(container[0]), nil

	default:
		return nil, fmt.Errorf("invalid type given to '[': %T", container)
	}
}

// tail returns a list/string of everything but the first element/char.
//
// This returns an error if the argument is not a list or string, or is empty.
func tail(args []Value) (Value, error) {
	container, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch container := container.(type) {
	case List:
		if len(container) == 0 {
			return nil, errors.New("empty list given to ']'")
		}

		return container[1:], nil

	case String:
		if len(container) == 0 {
			return nil, errors.New("empty string given to ']'")
		}

		return container[1:], nil

	default:
		return nil, fmt.Errorf("invalid type given to ']': %T", container)
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
	exitStatus, err := runTo[Integer](args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(int(exitStatus))
	panic("<unreachable>")
}

// not returns the logical negation of its argument
func not(args []Value) (Value, error) {
	boolean, err := runTo[Boolean](args[0])
	if err != nil {
		return nil, err
	}

	return Boolean(!boolean), nil
}

// length converts its argument to a list, then returns its length.
//
// This returns an error if the argument is not a list or string.
func length(args []Value) (Value, error) {
	container, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	// There need to be two conditions, even though they're identical, as the `len` function is
	// operating on a different type (and so `case List, String:` wouldn't work).
	switch container := container.(type) {
	case List:
		return Integer(len(container)), nil

	case String:
		return Integer(len(container)), nil

	default:
		list, err := TryConvert[List](container)
		if err != nil {
			return nil, err
		}
		return Integer(len(list)), nil

		// default:
		// 	return nil, fmt.Errorf("invalid type given to 'LENGTH': %T", container)
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
	string, err := runTo[String](args[0])
	if err != nil {
		return nil, err
	}

	if string != "" && string[len(string)-1] == '\\' {
		fmt.Print(string[:len(string)-1])
	} else {
		fmt.Println(string)
	}

	return Null{}, nil
}

// ascii is essentially equivalent to `chr`/`ord` in other langauges, depending on its argument.
func ascii(args []Value) (Value, error) {
	value, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch value := value.(type) {
	case Integer:
		if !utf8.ValidRune(rune(value)) {
			return nil, fmt.Errorf("invalid integer given to 'ASCII': %d", value)
		}

		return String(rune(value)), nil

	case String:
		if value == "" {
			return nil, errors.New("empty string given to 'ASCII'")
		}

		rune, _ := utf8.DecodeRuneInString(string(value))
		return Integer(rune), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'ASCII': %T", value)
	}
}

// negate returns the numerical negation of its argument.
func negate(args []Value) (Value, error) {
	number, err := runTo[Integer](args[0])
	if err != nil {
		return nil, err
	}

	return -number, nil
}

/** ARITY TWO **/

// add adds two numbers/strings/lists together; it coerces the second argument.
func add(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := runTo[Integer](args[1])
		if err != nil {
			return nil, err
		}

		return lhs + rhs, nil

	case String:
		rhs, err := runTo[String](args[1])
		if err != nil {
			return nil, err
		}

		// using `strings.Builder` is a bit more efficient than concating and stuff.
		var sb strings.Builder
		sb.WriteString(string(lhs))
		sb.WriteString(string(rhs))

		return String(sb.String()), nil

	case List:
		rhs, err := runTo[List](args[1])
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
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := runTo[Integer](args[1])
		if err != nil {
			return nil, err
		}

		return lhs - rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '-': %T", lhs)
	}
}

// multiply multiplies two numbers, or repeats lists/strings; last argument's converted to a number.
func multiply(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	// It just so happens that all three multiply cases need integers as the second argument
	rhs, err := runTo[Integer](args[1])
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		return lhs * rhs, nil

	case String:
		if rhs < 0 {
			return nil, fmt.Errorf("negative replication amount for a string in '*': %d", rhs)
		}

		return String(strings.Repeat(string(lhs), int(rhs))), nil

	case List:
		if rhs < 0 {
			return nil, fmt.Errorf("negative replication amount for a list in '*': %d", rhs)
		}

		slice := make(List, 0, len(lhs)*int(rhs))
		for i := 0; i < int(rhs); i++ {
			slice = append(slice, lhs...)
		}

		return slice, nil

	default:
		return nil, fmt.Errorf("invalid type given to '*': %T", lhs)
	}
}

// divide divides the first argument by the second; errors out of the second is zero.
func divide(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := runTo[Integer](args[1])
		if err != nil {
			return nil, err
		}

		if rhs == 0 {
			return nil, errors.New("zero divisor given to '/'")
		}

		return lhs / rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '/': %T", lhs)
	}
}

// remainder returns the remainder of `<arg1>/<arg2>`; errors out if second arg is zero.
func remainder(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := runTo[Integer](args[1])
		if err != nil {
			return nil, err
		}

		if rhs == 0 {
			return nil, errors.New("zero divisor given to '%'")
		}

		return lhs % rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '%%': %T", lhs)
	}
}

// exponentiate raises the first argument to the power of the second, or joins lists. errors out on
// negative powers for integers.
func exponentiate(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := runTo[Integer](args[1])
		if err != nil {
			return nil, err
		}

		if rhs < 0 {
			return nil, fmt.Errorf("negative exponent given to '^': %d", rhs)
		}

		// All 32 bit number exponentiations that can be represented in 32 bits can be done with
		// 64 bit floats and a "powf" function.
		return Integer(math.Pow(float64(lhs), float64(rhs))), nil

	case List:
		sep, err := runTo[String](args[1])
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
		rhs, err := TryConvert[Integer](rhs)
		if err != nil {
			return 0, err
		}

		return int(lhs - rhs), nil

	case String:
		rhs, err := TryConvert[String](rhs)
		if err != nil {
			return 0, err
		}

		return strings.Compare(string(lhs), string(rhs)), nil

	case Boolean:
		rhs, err := TryConvert[Boolean](rhs)
		if err != nil {
			return 0, err
		}

		return int(lhs.ToInteger() - rhs.ToInteger()), nil

	case List:
		rhs, err := TryConvert[List](rhs)
		if err != nil {
			return 0, err
		}

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
func lessThan(args []Value) (Value, error) {
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
func greaterThan(args []Value) (Value, error) {
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
func equalTo(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	rval, err := args[1].Run()
	if err != nil {
		return nil, err
	}

	// `DeepEqual` happens to correspond exactly to Knight's equality semantics
	return Boolean(reflect.DeepEqual(lhs, rval)), nil
}

// and evaluates the first argument, then either returns that if it's truthy, or otherwise evaluates
// and returns the second argument.
func and(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	isTruthy, err := TryConvert[Boolean](lhs)
	if err != nil {
		return nil, err
	}

	if isTruthy {
		return args[1].Run()
	}

	return lhs, nil
}

// or evaluates the first argument, then either returns that if it's falsey, or otherwise evaluates
// and returns the second argument.
func or(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	isTruthy, err := TryConvert[Boolean](lhs)
	if err != nil {
		return nil, err
	}

	if !isTruthy {
		return args[1].Run()
	}
	return lhs, nil
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
		condition, err := runTo[Boolean](args[0])
		if err != nil {
			return nil, err
		}

		if !condition {
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
	condition, err := runTo[Boolean](args[0])
	if err != nil {
		return nil, err
	}

	if condition {
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

	start, err := runTo[Integer](args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'GET': %d", start)
	}

	length, err := runTo[Integer](args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to GET: '%d'", length)
	}

	stop := start + length

	switch collection := collection.(type) {
	case String:
		if len(collection) < int(stop) {
			return nil, fmt.Errorf("string index out of bounds for 'GET': %d < %d", len(collection), stop)
		}

		return collection[start:stop], nil

	case List:
		if len(collection) < int(stop) {
			return nil, fmt.Errorf("list index out of bounds for 'GET': %d < %d", len(collection), stop)
		}

		return collection[start:stop], nil

	default:
		return nil, fmt.Errorf("invalid type given to 'GET': %T", collection)
	}
}

/** ARITY FOUR **/

// set returns a list/string where the range `[<arg2>, <arg2>+<arg3>)` is replaced by the fourth.
func set(args []Value) (Value, error) {
	collection, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	start, err := runTo[Integer](args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'SET': %d", start)
	}

	length, err := runTo[Integer](args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to 'SET': %d", length)
	}

	stop := start + length

	switch collection := collection.(type) {
	case String:
		if len(collection) < int(stop) {
			return nil, fmt.Errorf("string index out of bounds for 'SET': %d < %d", len(collection), stop)
		}

		replacement, err := runTo[String](args[3])
		if err != nil {
			return nil, err
		}

		var sb strings.Builder
		sb.WriteString(string(collection[:start]))
		sb.WriteString(string(replacement))
		sb.WriteString(string(collection[stop:]))
		return String(sb.String()), nil

	case List:
		if len(collection) < int(stop) {
			return nil, fmt.Errorf("list index out of bounds for 'SET': %d < %d", len(collection), stop)
		}

		replacement, err := runTo[List](args[3])
		if err != nil {
			return nil, err
		}

		return slices.Concat(collection[:start], replacement, collection[stop:]), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'SET': %T", collection)
	}
}
