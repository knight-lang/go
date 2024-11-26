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

// Function represents a Knight function (eg `DUMP`, `+`, `=`, etc.).
type Function struct {
	// The user-friendly name of the function. Used within syntax error and Ast.Dump.
	name string

	// The amount of arguments that `fn` expects.
	arity int

	// The go function associated with this function.
	fn func([]Value) (Value, error)
}

var (
	// KnownFunctions is a list of all known functions. The Parser uses this to recognize functions
	// in the source code, so modifying this map will change what functions the Parser knows about.
	KnownFunctions = map[rune]*Function{
		// Arity 0
		'T': &Function{name: "TRUE", arity: 0, fn: true_},
		'F': &Function{name: "FALSE", arity: 0, fn: false_},
		'N': &Function{name: "NULL", arity: 0, fn: null},
		'@': &Function{name: "@", arity: 0, fn: emptyList},
		'P': &Function{name: "PROMPT", arity: 0, fn: prompt},
		'R': &Function{name: "RANDOM", arity: 0, fn: random},

		// Arity 1
		':': &Function{name: ":", arity: 1, fn: noop},
		'B': &Function{name: "BLOCK", arity: 1, fn: block},
		'C': &Function{name: "CALL", arity: 1, fn: call},
		'Q': &Function{name: "QUIT", arity: 1, fn: quit},
		'!': &Function{name: "!", arity: 1, fn: not},
		'L': &Function{name: "LENGTH", arity: 1, fn: length},
		'D': &Function{name: "DUMP", arity: 1, fn: dump},
		'O': &Function{name: "OUTPUT", arity: 1, fn: output},
		'A': &Function{name: "ASCII", arity: 1, fn: ascii},
		'~': &Function{name: "~", arity: 1, fn: negate},
		',': &Function{name: ",", arity: 1, fn: box},
		'[': &Function{name: "[", arity: 1, fn: head},
		']': &Function{name: "]", arity: 1, fn: tail},

		// Arity 2
		'+': &Function{name: "+", arity: 2, fn: add},
		'-': &Function{name: "-", arity: 2, fn: subtract},
		'*': &Function{name: "*", arity: 2, fn: multiply},
		'/': &Function{name: "/", arity: 2, fn: divide},
		'%': &Function{name: "%", arity: 2, fn: remainder},
		'^': &Function{name: "^", arity: 2, fn: exponentiate},
		'<': &Function{name: "<", arity: 2, fn: lessThan},
		'>': &Function{name: ">", arity: 2, fn: greaterThan},
		'?': &Function{name: "?", arity: 2, fn: equalTo},
		'&': &Function{name: "&", arity: 2, fn: and},
		'|': &Function{name: "|", arity: 2, fn: or},
		';': &Function{name: ";", arity: 2, fn: then},
		'=': &Function{name: "=", arity: 2, fn: assign},
		'W': &Function{name: "WHILE", arity: 2, fn: while},

		// Arity 3
		'I': &Function{name: "IF", arity: 3, fn: if_},
		'G': &Function{name: "GET", arity: 3, fn: get},

		// Arity 4
		'S': &Function{name: "SET", arity: 4, fn: set},
	}

	// stdinScanner is used by prompt to read lines from the standard input.
	stdinScanner = bufio.NewScanner(os.Stdin)
)

// initialize the random number generator. (For non-go-folks, go ensures that all functions named
// `init` will be executed before `main` is run.)
func init() {
	rand.Seed(time.Now().UnixNano())
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 0                                             *
 *                                                                                                *
 **************************************************************************************************/

// "Literal" functions---functions which take no arguments and always return the same value.
func true_(_ []Value) (Value, error)     { return Boolean(true), nil }
func false_(_ []Value) (Value, error)    { return Boolean(false), nil }
func null(_ []Value) (Value, error)      { return Null{}, nil }
func emptyList(_ []Value) (Value, error) { return List{}, nil }

// random returns a random `Integer`.
func random(_ []Value) (Value, error) {
	// Note that `rand` is seeded in this file's `init` function.
	return Integer(rand.Int63()), nil
}

// prompt reads a line from stdin, returning Null if stdin is empty.
func prompt(_ []Value) (Value, error) {
	if stdinScanner.Scan() {
		return String(strings.TrimRight(stdinScanner.Text(), "\r")), nil
	}

	if err := stdinScanner.Err(); err != nil {
		return nil, err
	}

	return Null{}, nil
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 1                                             *
 *                                                                                                *
 **************************************************************************************************/

// noop simply executes its only argument and returns it
func noop(args []Value) (Value, error) {
	return args[0].Run()
}

// box creates a list just containing its argument.
func box(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return List{ran}, nil
}

// head returns the first element/rune of a list/string. It returns an error if the container is
// empty, or if the argument isn't a list or string.
func head(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch container := ran.(type) {
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

// tail returns a list/string of everything but the first element/rune. It returns an error if the
// container is empty, or if the argument isn't a list or string.
func tail(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch container := ran.(type) {
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

// block returns its argument unevaluated. This is intended to be used in conjunction with call (see
// below) to defer evaluation to a later point in time.
func block(args []Value) (Value, error) {
	return args[0], nil
}

// call runs its argument, and then returns the result of running _that_ value. This allows us to
// defer execution of `BLOCK`s until later on.
func call(args []Value) (Value, error) {
	block, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	return block.Run()
}

// quit exits the program with the given exit status code.
func quit(args []Value) (Value, error) {
	exitStatus, err := runToInteger(args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(int(exitStatus))
	panic("<unreachable>") // Go isn't smart enough to recognize `os.Exit` never returns.
}

// not returns the logical negation of its argument
func not(args []Value) (Value, error) {
	boolean, err := runToBoolean(args[0])
	if err != nil {
		return nil, err
	}

	return Boolean(!boolean), nil
}

// negate returns the numerical negation of its argument.
func negate(args []Value) (Value, error) {
	integer, err := runToInteger(args[0])
	if err != nil {
		return nil, err
	}

	return -integer, nil
}

// length returns the length of a list/string. It returns an error if the argument isn't a
// list or string.
func length(args []Value) (Value, error) {
	container, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	// (Note: There need to be two branches here even though their contents are identical because the
	// `len` function is operating on two different types).
	switch container := container.(type) {
	case List:
		return Integer(len(container)), nil

	case String:
		return Integer(len(container)), nil

	default:
		// Knight 2.0.1 required `LENGTH` to coerce its arguments to lists, instead of having it be
		// undefined behaviour
		if shouldSupportKnightVersion_2_0_1 {
			list, err := container.ToList()
			if err != nil {
				return nil, err
			}
			return Integer(len(list)), nil
		}

		return nil, fmt.Errorf("invalid type given to 'LENGTH': %T", container)

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

// output writes its argument to stdout and returns null. Normally, a newline is added after its
// argument, however if the argument ends in a `\`, the backslash is removed and no newline is
// printed.
func output(args []Value) (Value, error) {
	string, err := runToString(args[0])
	if err != nil {
		return nil, err
	}

	if string != "" && string[len(string)-1] == '\\' {
		fmt.Print(string[:len(string)-1])

		// Since we're not printing a newline, we flush stdout so that the output is always visible.
		// (The error is explicitly ignored to be consistent with how `fmt.Print{,ln}` works.)
		_ = os.Stdout.Sync()
	} else {
		fmt.Println(string)
	}

	return Null{}, nil
}

// ascii is the equivalent of `chr()` and `ord()` functions in other languages. An error is returned
// if an empty string, an integer which doesn't correspond to a rune, or a non int-non-string type
// is given.
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

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 2                                             *
 *                                                                                                *
 **************************************************************************************************/

// add adds two integers/strings/lists together by coercing the second argument. Passing in any
// other type will yield an error.
func add(args []Value) (Value, error) {
	ran, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := ran.(type) {
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

		return slices.Concat(lhs, rhs), nil

	default:
		return nil, fmt.Errorf("invalid type given to '+': %T", lhs)
	}
}

// subtract subtracts one integer from another. It returns an error for other types.
func subtract(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := runToInteger(args[1])
		if err != nil {
			return nil, err
		}

		return lhs - rhs, nil

	default:
		return nil, fmt.Errorf("invalid type given to '-': %T", lhs)
	}
}

// multiply an integer by another, or repeats a list or string. It returns an error for other types.
func multiply(args []Value) (Value, error) {
	lhs, err := args[0].Run()
	if err != nil {
		return nil, err
	}

	// It just so happens that all three multiply cases need integers as the second argument, so
	// just do the coercion before the typecheck.
	rhs, err := runToInteger(args[1])
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

		return slices.Repeat(lhs, int(rhs)), nil

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
		rhs, err := runToInteger(args[1])
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
		rhs, err := runToInteger(args[1])
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
		rhs, err := runToInteger(args[1])
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
		sep, err := runToString(args[1])
		if err != nil {
			return nil, err
		}

		joined, err := lhs.Join(string(sep))
		if err != nil {
			return nil, err
		}

		return String(joined), nil

	default:
		return nil, fmt.Errorf("invalid type given to '^': %T", lhs)
	}
}

func compare(lhs, rhs Value, fn rune) (int, error) {
	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := rhs.ToInteger()
		if err != nil {
			return 0, err
		}

		return int(lhs - rhs), nil

	case String:
		rhs, err := rhs.ToString()
		if err != nil {
			return 0, err
		}

		return strings.Compare(string(lhs), string(rhs)), nil

	case Boolean:
		rhs, err := rhs.ToBoolean()
		if err != nil {
			return 0, err
		}

		if lhs == rhs {
			return 0, nil
		} else if lhs && !rhs {
			return 1, nil
		} else {
			return -1, nil
		}

	case List:
		rhs, err := rhs.ToList()
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

	isTruthy, err := lhs.ToBoolean()
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

	isTruthy, err := lhs.ToBoolean()
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
		condition, err := runToBoolean(args[0])
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
	condition, err := runToBoolean(args[0])
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

	start, err := runToInteger(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'GET': %d", start)
	}

	length, err := runToInteger(args[2])
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

	start, err := runToInteger(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'SET': %d", start)
	}

	length, err := runToInteger(args[2])
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

		replacement, err := runToString(args[3])
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

		replacement, err := runToList(args[3])
		if err != nil {
			return nil, err
		}

		return slices.Concat(collection[:start], replacement, collection[stop:]), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'SET': %T", collection)
	}
}
