package knight

import (
	"bufio"
	"errors" // For those non-gophers, `errors.New` is `fmt.Errorf` when no interpolation is needed.
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// Function represents a Knight function (eg `DUMP`, `+`, `=`, etc.).
//
// These are used within Ast to store which function the AST should be executing.
type Function struct {
	// The user-friendly name of the function. Used within syntax error and `Ast.Dump`.
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

	// stdinScanner is used by the `prompt` function to read lines from the standard input.
	stdinScanner = bufio.NewScanner(os.Stdin)
)

// Initialize the functions module. This both initializes the random number generator for `random`,
// as well as registers extension functions.
//
// (For non-go-folks, go ensures that each file's `init` function, if it exists, will be executed
// before `main` is run.)
func init() {
	rand.Seed(time.Now().UnixNano())

	// Extension functions
	KnownFunctions['E'] = &Function{name: "EVAL", arity: 1, fn: eval}
	KnownFunctions['$'] = &Function{name: "$", arity: 1, fn: system}
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 0                                             *
 *                                                                                                *
 **************************************************************************************************/

// true_ always returns the true Boolean.
func true_(_ []Value) (Value, error) {
	return Boolean(true), nil
}

// false_ always returns the false Boolean.
func false_(_ []Value) (Value, error) {
	return Boolean(false), nil
}

// null always returns Null.
func null(_ []Value) (Value, error) {
	return Null{}, nil
}

// emptyList always returns an empty List
func emptyList(_ []Value) (Value, error) {
	return List{}, nil
}

// random returns a random Integer.
func random(_ []Value) (Value, error) {
	// Note that `rand` is seeded in this file's `init` function.
	return Integer(rand.Int63()), nil // Go only has `Int63` for some reason...
}

// prompt reads a line from stdin, returning Null if stdin is empty.
func prompt(_ []Value) (Value, error) {
	// If there was a problem getting the line, then we're either at the end of the file (which means
	// we should return Null), or there was some problem like stdin was closed or permission denied.
	if !stdinScanner.Scan() {
		// EOF doesn't cause errors; this means there's a problem with stdin, like permission denied.
		if err := stdinScanner.Err(); err != nil {
			return nil, fmt.Errorf("unable to 'PROMPT': %v", err)
		}

		// EOF was reached, return null.
		return Null{}, nil
	}

	// The line was scanned properly, extract it.
	line := stdinScanner.Text()

	// Knight version 2.0.1 requires _all_ trailing `\r`s to be stripped. Knight 2.1 only requires
	// one to be stripped (which `.Text()` does for us).
	if shouldSupportKnightVersion_2_0_1 {
		line = strings.TrimRight(line, "\r")
	}

	return String(line), nil
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 1                                             *
 *                                                                                                *
 **************************************************************************************************/

// noop simply executes its only argument and returns it
func noop(args []Value) (Value, error) {
	return args[0].Execute()
}

// box creates a list just containing its argument.
func box(args []Value) (Value, error) {
	ran, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	return List{ran}, nil
}

// head returns the first element/rune of a list/string. It returns an error if the container is
// empty, or if the argument isn't a list or string.
func head(args []Value) (Value, error) {
	ran, err := args[0].Execute()
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
	ran, err := args[0].Execute()
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

// block returns its argument unexecuted. This is intended to be used in conjunction with call (see
// below) to defer evaluation to a later point in time.
func block(args []Value) (Value, error) {
	return args[0], nil
}

// call executes its argument, and then returns the result of executing _that_ value. This allows us
// to defer execution of `BLOCK`s until later on.
func call(args []Value) (Value, error) {
	block, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	return block.Execute()
}

// quit exits the program with the given exit status code.
func quit(args []Value) (Value, error) {
	exitStatus, err := executeToInteger(args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(int(exitStatus))
	panic("<unreachable>") // Go isn't powerful enough to recognize os.Exit never returns.
}

// not returns the logical negation of its argument
func not(args []Value) (Value, error) {
	boolean, err := executeToBoolean(args[0])
	if err != nil {
		return nil, err
	}

	return !boolean, nil
}

// negate returns the numerical negation of its argument.
func negate(args []Value) (Value, error) {
	integer, err := executeToInteger(args[0])
	if err != nil {
		return nil, err
	}

	return -integer, nil
}

// length returns the length of a list/string. It returns an error if the argument isn't a
// list or string.
func length(args []Value) (Value, error) {
	container, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	// (Note: There need to be two branches here even though their contents are identical because the
	// `len` function is operating on two different types: In the first case, a List, in the second
	// a String.)
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
	value, err := args[0].Execute()
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
	message, err := executeToString(args[0])
	if err != nil {
		return nil, err
	}

	// Get the last "rune" (go-speak for (ish) a unicode character), so we can compare it against a
	// backslash to see if the string ends in `\`. (If it does, the Knight specs say it should be
	// deleted and the normal newline that `OUTPUT` would print would be suppressed.)
	// NOTE: `DecodeLastRuneInString` will return `RuneError` if the message is empty. Since we only
	// compare it against backslash, we don't need explicitly check for `string`'s length.
	lastChr, idx := utf8.DecodeLastRuneInString(string(message))

	if lastChr == '\\' {
		fmt.Print(message[:len(message)-idx])

		// Since we're not printing a newline, we flush stdout so that the output is always visible.
		// (The error is explicitly ignored to be consistent with how `fmt.Print{,ln}` works.)
		_ = os.Stdout.Sync()
	} else {
		fmt.Println(message)
	}

	return Null{}, nil
}

// ascii is the equivalent of `chr()` and `ord()` functions in other languages. An error is returned
// if an empty string, an integer which doesn't correspond to a rune, or a non int-non-string type
// is given.
func ascii(args []Value) (Value, error) {
	value, err := args[0].Execute()
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
	ran, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := ran.(type) {
	case Integer:
		rhs, err := executeToInteger(args[1])
		if err != nil {
			return nil, err
		}

		return lhs + rhs, nil

	case String:
		rhs, err := executeToString(args[1])
		if err != nil {
			return nil, err
		}

		// using strings.Builder is a bit more efficient than concating and stuff.
		var sb strings.Builder
		sb.WriteString(string(lhs))
		sb.WriteString(string(rhs))
		return String(sb.String()), nil

	case List:
		rhs, err := executeToList(args[1])
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
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInteger(args[1])
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
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	// It just so happens that all three multiply cases need integers as the second argument, so
	// just do the coercion before the typecheck.
	rhs, err := executeToInteger(args[1])
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

// divide divides an integer by another. It returns an error for other types, or if the second
// argument is zero.
func divide(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInteger(args[1])
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

// remainder gets the remainder of the first argument and the second. It returns an error for other
// types, or if the second argument is zero.
func remainder(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInteger(args[1])
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

// exponentiate raises the first argument to the power of the second, or joins lists. It returns an
// error for other types, if an integer is raised to a negative power, or if the list contains types
// which cannot be converted to strings (such as `BLOCK`'s return value).
func exponentiate(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInteger(args[1])
		if err != nil {
			return nil, err
		}

		if rhs < 0 {
			return nil, fmt.Errorf("negative exponent given to '^': %d", rhs)
		}

		// Knight only requires us support 32 bit integers, and only support exponentiations which
		// don't overflow those bounds. This requirement can be satisfied by converting to float64s,
		// as they can losslessly represent 32 bit integers. While this does mean that excessively
		// large 64 bit integers won't yield exactly correct results, this method is much faster and
		// cleaner than having to do exponentiation ourselves.
		return Integer(math.Pow(float64(lhs), float64(rhs))), nil

	case List:
		sep, err := executeToString(args[1])
		if err != nil {
			return nil, err
		}

		joined, err := lhs.Join(string(sep)) // Join can fail if the list contains Asts or Variables.
		if err != nil {
			return nil, err
		}

		return String(joined), nil

	default:
		return nil, fmt.Errorf("invalid type given to '^': %T", lhs)
	}
}

// compare returns a negative, zero, or positive integer depending on whether lhs is less than,
// equal to, or greater than the second. The functionName argument is just used for error messages
// if an invalid type is provided.
func compare(lhs, rhs Value, functionName rune) (int, error) {
	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := rhs.ToInteger()
		if err != nil {
			return 0, err
		}

		// Subtraction actually is all that's needed for integers.
		return int(lhs - rhs), nil

	case String:
		rhs, err := rhs.ToString()
		if err != nil {
			return 0, err
		}

		// strings.Compare does lexicographical comparisons
		return strings.Compare(string(lhs), string(rhs)), nil

	case Boolean:
		rhs, err := rhs.ToBoolean()
		if err != nil {
			return 0, err
		}

		// Just manually enumerate all the cases for booleans.
		if !lhs && rhs {
			return -1, nil
		} else if lhs && !rhs {
			return 1, nil
		} else {
			return 0, nil
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

		// Check element-wise, and return the first non-equal comparison.
		for i := 0; i < minLen; i++ {
			cmp, err := compare(lhs[i], rhs[i], functionName)
			if err != nil {
				return 0, err
			}

			if cmp != 0 {
				return cmp, nil
			}
		}

		// All elements were equal, now check their lengths.
		return len(lhs) - len(rhs), nil

	default:
		return 0, fmt.Errorf("invalid type given to %q: %T", functionName, lhs)
	}
}

// lessThan returns whether the first argument is less than the second. An error is returned if
// the first argument isn't a boolean, integer, string, or list, or if a list that's passed contains
// an invalid argument.
func lessThan(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	rhs, err := args[1].Execute()
	if err != nil {
		return nil, err
	}

	cmp, err := compare(lhs, rhs, '<')
	if err != nil {
		return nil, err
	}

	return Boolean(cmp < 0), nil
}

// greaterThan returns whether the first argument is greater than the second. An error is returned
// if the first argument isn't a boolean, integer, string, or list, or if a list that's passed
// contains an invalid argument.
func greaterThan(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	rhs, err := args[1].Execute()
	if err != nil {
		return nil, err
	}

	cmp, err := compare(lhs, rhs, '>')
	if err != nil {
		return nil, err
	}

	return Boolean(cmp > 0), nil
}

// equalTo returns whether its two arguments are equal to one other. Unlike the `<` and `>`
// functions, this doesn't coerce the second argument to the type of the first.
func equalTo(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	rval, err := args[1].Execute()
	if err != nil {
		return nil, err
	}

	// reflect.DeepEqual happens to correspond exactly to Knight's equality semantics.
	return Boolean(reflect.DeepEqual(lhs, rval)), nil
}

// and evaluates the first argument and returns it if it's falsey. When it's truthy, it returns the
// second argument.
func and(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	isTruthy, err := lhs.ToBoolean()
	if err != nil {
		return nil, err
	}

	if isTruthy {
		return args[1].Execute()
	}

	return lhs, nil
}

// or evaluates the first argument and returns it if it's truthy. When it's falsey, it returns the
// second argument.
func or(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	isTruthy, err := lhs.ToBoolean()
	if err != nil {
		return nil, err
	}

	if !isTruthy {
		return args[1].Execute()
	}
	return lhs, nil
}

// then evaluates the first argument, then evaluates and returns the second argument.
func then(args []Value) (Value, error) {
	if _, err := args[0].Execute(); err != nil {
		return nil, err
	}

	return args[1].Execute()
}

// assign is used to assign values to variables. The first argument must be a Variable, or an error
// is returned. The second argument is evaluated, and after assignment is returned.
func assign(args []Value) (Value, error) {
	variable, ok := args[0].(*Variable)
	if !ok {
		return nil, fmt.Errorf("invalid type given to '=': %T", args[0])
	}

	value, err := args[1].Execute()
	if err != nil {
		return nil, err
	}

	variable.Assign(value)

	return value, nil
}

// while evaluates the second argument whilst the first is true, and returns Null.
func while(args []Value) (Value, error) {
	for {
		condition, err := executeToBoolean(args[0])
		if err != nil {
			return nil, err
		}

		if !condition {
			break
		}

		if _, err = args[1].Execute(); err != nil {
			return nil, err
		}
	}

	return Null{}, nil
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 2                                             *
 *                                                                                                *
 **************************************************************************************************/

// if_ evaluates and returns the second argument if the first is truthy; if it's falsey, if_
// evaluates and returns the third argument instead.
func if_(args []Value) (Value, error) {
	condition, err := executeToBoolean(args[0])
	if err != nil {
		return nil, err
	}

	if condition {
		return args[1].Execute()
	}

	return args[2].Execute()
}

// get returns a sublist/substring with start and length of the second and third arguments. It
// returns an error if the start or length are negative, if `start + length` is larger than
// the collection's length, or if a non-list/string element is provided.
func get(args []Value) (Value, error) {
	collection, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	start, err := executeToInteger(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'GET': %d", start)
	}

	length, err := executeToInteger(args[2])
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

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 4                                             *
 *                                                                                                *
 **************************************************************************************************/

// set returns a list/string where the range `[start, start+length)` (where start and length are the
// second and third parameters, respectively) is replaced by the fourth parameter. An error is
// returned if either the start or length are negative, if `start+length` is larger than the size
// of the container, or if the first argument isn't a list or string.
func set(args []Value) (Value, error) {
	collection, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	start, err := executeToInteger(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'SET': %d", start)
	}

	length, err := executeToInteger(args[2])
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

		replacement, err := executeToString(args[3])
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

		replacement, err := executeToList(args[3])
		if err != nil {
			return nil, err
		}

		return slices.Concat(collection[:start], replacement, collection[stop:]), nil

	default:
		return nil, fmt.Errorf("invalid type given to 'SET': %T", collection)
	}
}

/**************************************************************************************************
 *                                                                                                *
 *                                           Extensions                                           *
 *                                                                                                *
 **************************************************************************************************/

func eval(args []Value) (Value, error) {
	sourceCode, err := executeToString(args[0])
	if err != nil {
		return nil, err
	}

	return Evaluate(string(sourceCode))
}

func system(args []Value) (Value, error) {
	// Get the shell script to execute
	shellCommand, err := executeToString(args[0])
	if err != nil {
		return nil, err
	}

	// Use the `SHELL` environment variable, if it exists. If it doesn't, default to `/bin/sh`
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	// Execute the command
	stdout, err := exec.Command(shell, "-c", string(shellCommand)).Output()
	if err != nil {
		return nil, err
	}

	// Return the stdout
	return String(stdout), nil
}
