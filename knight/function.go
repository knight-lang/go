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
// These are used within FnCall to store which function the function call should be executing.
type Function struct {
	// User-friendly name of the function. Used within syntax error and `FnCall.Dump`.
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

	// Extension functions. (We have to add these here because including `eval` above would be a
	// circular loop; I moved `system` out here to be consistent)
	KnownFunctions['E'] = &Function{name: "EVAL", arity: 1, fn: eval}
	KnownFunctions['`'] = &Function{name: "`", arity: 1, fn: system}
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 0                                             *
 *                                                                                                *
 **************************************************************************************************/

// true_ always returns the true Boolean.
//
// ## Examples
//
//	DUMP TRUE #=> true
func true_(_ []Value) (Value, error) {
	return Boolean(true), nil
}

// false_ always returns the false Boolean.
//
// ## Examples
//
//	DUMP FALSE #=> false
func false_(_ []Value) (Value, error) {
	return Boolean(false), nil
}

// null always returns Null.
//
// ## Examples
//
//	DUMP NULL #=> null
func null(_ []Value) (Value, error) {
	return Null{}, nil
}

// emptyList always returns an empty List
//
// ## Examples
//
//	DUMP @ #=> []
func emptyList(_ []Value) (Value, error) {
	return List{}, nil
}

// random returns a random Integer.
//
// As an extension, the go implementation supports random integers above the required 32767.
//
// ## Examples
//
//	DUMP RANDOM #=> 8015671084101644486
func random(_ []Value) (Value, error) {
	// Note that `rand` is seeded in this file's `init` function.
	return Integer(rand.Int63()), nil // Go only has `Int63` for some reason...
}

// prompt reads a line from stdin, returning Null if stdin is empty.
//
// ## Examples
//
//	DUMP PROMPT <stdin="foo">        #=> "foo"
//	DUMP PROMPT <stdin="foo\n">      #=> "foo"
//	DUMP PROMPT <stdin="foo\nbar">   #=> "foo"
//	DUMP PROMPT <stdin="foo\r\nbar"> #=> "foo"
//	DUMP PROMPT <stdin="foo\rbar">   #=> "foo\rbar"
//	DUMP PROMPT <stdin="foo\r">      #=> "foo"
//	DUMP PROMPT <stdin="">           #=> ""
//	DUMP ; PROMPT PROMPT <stdin="">  #=> null
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

	// The line was scanned properly, return it.
	return String(stdinScanner.Text()), nil
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 1                                             *
 *                                                                                                *
 **************************************************************************************************/

// noop simply executes its only argument and returns it
//
// ## Examples
//
//	DUMP : 34                     #=> 34
//	: : : DUMP : : : : + : 30 : 4 #=> 34
//
// : (BLOCK foo)                 # (works, `:` accepts Blocks)
func noop(args []Value) (Value, error) {
	return args[0].Execute()
}

// box creates a list just containing its argument.
//
// ## Examples
//
//	DUMP ,T        #=> [true]
//	DUMP ,@        #=> [[]]
//	DUMP ,,,,3     #=> [[[[3]]]]
//
// , (BLOCK foo)  # (works, `,` accepts Blocks)
func box(args []Value) (Value, error) {
	ran, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	return List{ran}, nil
}

// head returns the first element/rune of a list/string. It returns an error if the container is
// empty, or if the argument isn't a list or string.
//
// ## Examples
//
//	DUMP [ "A"   #=> "A"
//	DUMP [ "ABC" #=> "A"
//	DUMP [ ,1    #=> 1
//	DUMP [ +@123 #=> 1
//
// ## Undefined Behaviour
// Errors are returned for all forms of undefined behaviour in `[`:
//
//	DUMP [ ""    #!! empty string
//	DUMP [ @     #!! empty list
//	DUMP [ 123   #!! other types
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
//
// ## Examples
//
//	DUMP ] "A"   #=> ""
//	DUMP ] "ABC" #=> "BC"
//	DUMP ] ,1    #=> []
//	DUMP ] +@123 #=> [2, 3]
//
// ## Undefined Behaviour
// Errors are returned for all forms of undefined behaviour in `]`:
//
//	DUMP ] ""    #!! empty string
//	DUMP ] @     #!! empty list
//	DUMP ] 123   #!! other types
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
//
// Because the argument is returned, unexecuted, we might be returning a `Variable` or an `FnCall`,
// both of which have no conversions defined on them. However, since using the return value of
// `BLOCK` in any function other than `CALL` (or, a handful of others that don't actually modify,
// their argument, such as `;`, `&`'s second argument, etc.) is undefined behaviour, this is a
// totally legit and fine strategy.
//
// ## Examples
//
//	; = double BLOCK * x 2
//	; = x 2
//	; OUTPUT CALL double     #=> 4
//	; = x 10
//	: OUTPUT CALL double     #=> 20
//
// BLOCK (BLOCK foo)        # (works, `BLOCK` also accepts Blocks)
func block(args []Value) (Value, error) {
	return args[0], nil
}

// call executes its argument, and then returns the result of executing _that_ value. This allows us
// to defer execution of `BLOCK`s until later on.
//
// ## Examples
//
//	; = double BLOCK * x 2
//	; = x 2
//	; OUTPUT CALL double     #=> 4
//	; = x 10
//	: OUTPUT CALL double     #=> 20
//
// ## Undefined Behaviour
// `CALL`ing non-`BLOCK`s is supported, and simply returns that argument:
//
//	DUMP CALL 123 #=> 123
//
// (NOTE: This is a direct consequence of how `BLOCK` is implemented, as `BLOCK 12` actually
// returns `12`, so `CALL BLOCK 12` actually reduces down to `CALL 12`, which then returns `12`.)
func call(args []Value) (Value, error) {
	block, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	return block.Execute()
}

// quit exits the program with the given exit status code.
//
// ## Examples
//
//	QUIT TRUE   # (exit with status 1)
//	QUIT NULL   # (exit with status 0)
//	QUIT 12     # (exit with status 12)
//	QUIT "017"  # (exit with status 17)
//
// ## Undefined Behaviour
// As an extension, exit codes that can fit into an `int` are supported. (Although, the OS might
// not let us return them.)
//
//	QUIT 12345  # (allowed, but the OS determines the exit status...)
func quit(args []Value) (Value, error) {
	exitStatus, err := executeToInt(args[0])
	if err != nil {
		return nil, err
	}

	os.Exit(exitStatus)
	panic("<unreachable>") // Go isn't powerful enough to recognize os.Exit never returns.
}

// not returns the logical negation of its argument
//
// ## Examples
//
//	DUMP ! NULL #=> true
//	DUMP ! @     #=> true
//	DUMP ! TRUE  #=> false
//	DUMP ! 12    #=> false
//	DUMP ! "0"   #=> false
//	DUMP ! ,@    #=> false
//
// ## Undefined Behaviour
// Types which can't be converted to booleans yield an error:
//
//	DUMP ! BLOCK foo    #!! error: cant convert to a boolean
func not(args []Value) (Value, error) {
	boolean, err := executeToBool(args[0])
	if err != nil {
		return nil, err
	}

	return Boolean(!boolean), nil
}

// negate returns the numerical negation of its argument.
//
// ## Examples
//
//	DUMP ~ FALSE #=> 0
//	DUMP ~ @     #=> 0
//	DUMP ~ TRUE  #=> -1
//	DUMP ~ 12    #=> -12
//	DUMP ~ "-12" #=> 12
//	DUMP ~ ~ 12  #=> 12
//	DUMP ~ "017" #=> -17
//	DUMP ~ "hi"  #=> 0
//	DUMP ~ ,@    #=> -1
//
// ## Undefined Behaviour
// Types which can't be converted to booleans yield an error:
//
//	DUMP ~ BLOCK foo    #!! error: cant convert to an integer
func negate(args []Value) (Value, error) {
	integer, err := executeToInt(args[0])
	if err != nil {
		return nil, err
	}

	return Integer(-integer), nil
}

// length returns the length of its argument, converted to an array.
//
// ## Examples
//
//	DUMP LENGTH 123            #=> 3
//	DUMP LENGTH "hello world"  #=> 11
//	DUMP LENGTH ++++,T,F,T,F,N #=> 5
//
// ## Undefined Behaviour
// Types which can't be converted to lists yield an error:
//
//	DUMP LENGTH BLOCK foo      #!! error: cant convert to a list
func length(args []Value) (Value, error) {
	list, err := executeToSlice(args[0])
	if err != nil {
		return nil, err
	}

	return Integer(len(list)), nil
}

// dump prints a debugging representation of its argument to stdout, then returns it.
//
// ## Examples
//
//	DUMP 123               #=> 123
//	DUMP 'he said\: "hi!"' #=> "he said\\: \"hi!\""
//	DUMP TRUE              #=> true
//
// ## Undefined Behaviour
// As an extension, _all_ types can be passed to `DUMP`.
//
//	DUMP BLOCK + foo 2     #=> FnCall(%!c(string=+), Variable(foo), 2)
//
// Any errors with writing to stdout are silently ignored.
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
//
// Examples (`␤` represents newline, to make these examples clearer):
//
//	OUTPUT 123             #=> 123␤
//	OUTPUT 'hello'         #=> hello␤
//	OUTPUT NULL            #=> ␤         (Note: `NULL` coerces to an empty string)
//	OUTPUT "what\"         #=> what      (notice there's no newline)
//	OUTPUT +@"a\"          #=> a␤        (trailng `\` was deleted, but the `␤` separator is kept)
//	OUTPUT "a\␤"           #=> a\␤␤
//	OUTPUT "a\␤\"          #=> a\␤
//
// ## Undefined Behaviour
// Types which can't be converted to strings yield an error:
//
//	OUTPUT BLOCK foo       #!! error: cant convert to a list
//
// Any errors with writing to stdout are silently ignored.
func output(args []Value) (Value, error) {
	message, err := executeToString(args[0])
	if err != nil {
		return nil, err
	}

	// Get the last "rune" (go-speak for (ish) a unicode character), so we can compare it against a
	// backslash to see if the string ends in `\`. (If it does, the Knight specs say it should be
	// deleted and the normal newline that `OUTPUT` would print would be suppressed.)
	//
	// NOTE: `DecodeLastRuneInString` will return `RuneError` if the message is empty. Since we only
	// compare it against backslash, we don't need explicitly check for `string`'s length.
	lastChr, idx := utf8.DecodeLastRuneInString(message)

	// Check to see if the last character is a `\`, and if it is, print neither it nor the newline
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
//
// ## Examples
//
//	DUMP ASCII 10    #=> <newline>
//	DUMP ASCII 126   #=> ~
//	DUMP ASCII "F"   #=> 70
//	DUMP ASCII "FOO" #=> 70
//
// ## Undefined Behaviour
// As an extension, `ASCII` supports all of utf-8:
//
//	DUMP ASCII "😁"   #=> 128513
//	DUMP ASCII 128513 #=> 😁
//
// Errors are returned for all other forms of undefined behaviour in `ASCII`:
//
//	DUMP ASCII 1234567 #!! not a valid rune
//	DUMP ASCII ""      #!! empty rune
//
// Other types are invalid:
//
//	DUMP ASCII TRUE    #!! error: invalid type
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
//
// ## Examples
//
//	DUMP + 12 34       #=> 36
//	DUMP + "hi" TRUE   #=> hitrue
//	DUMP + @ "what"    #=> ["w", "h", "a", "t"]
//
// ## Undefined Behaviour
// Overflowing operations on `Integer`s just wraparound
//
//	DUMP + 9223372036854775807 1                    #=> -9223372036854775808
//
// Creating lists or strings which are larger than `2147483647` will do whatever the golang runtime
// would do. (Which probably is a memory allocation error, and aborting the program.)
//
//	DUMP + "<2147483647-character-long string>" "X" #=> might work, depending on the OS
//
// Other types are invalid:
//
//	DUMP + TRUE 34  #!! error: invalid type
func add(args []Value) (Value, error) {
	ran, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := ran.(type) {
	case Integer:
		rhs, err := executeToInt(args[1])
		if err != nil {
			return nil, err
		}

		return Integer(int(lhs) + rhs), nil

	case String:
		rhs, err := executeToString(args[1])
		if err != nil {
			return nil, err
		}

		// using strings.Builder is a bit more efficient than concating and stuff.
		var sb strings.Builder
		sb.WriteString(string(lhs))
		sb.WriteString(rhs)
		return String(sb.String()), nil

	case List:
		rhs, err := executeToSlice(args[1])
		if err != nil {
			return nil, err
		}

		return slices.Concat(lhs, rhs), nil

	default:
		return nil, fmt.Errorf("invalid type given to '+': %T", lhs)
	}
}

// subtract subtracts one integer from another. It returns an error for other types.
//
// ## Examples
//
//	DUMP - 12 "34"       #=> -22
//	DUMP - 12 FALSE      #=> -12
//
// ## Undefined Behaviour
// Overflowing operations on `Integer`s just wraparound
//
//	DUMP - ~2 9223372036854775807 #=> 9223372036854775807
//
// Other types are invalid:
//
//	DUMP - TRUE 34  #!! error: invalid type
func subtract(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInt(args[1])
		if err != nil {
			return nil, err
		}

		return Integer(int(lhs) - rhs), nil

	default:
		return nil, fmt.Errorf("invalid type given to '-': %T", lhs)
	}
}

// multiply an integer by another, or repeats a list or string. It returns an error for other types.
//
// ## Examples
//
//	DUMP * 12 34     #=> 408
//	DUMP * "hi" 3    #=> hihihi
//	DUMP * (+@123) 4 #=> [1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3]
//
// ## Undefined Behaviour
// Overflowing operations on `Integer`s just wraparound
//
//	DUMP * 922337203685477580 20  #=> -16
//
// Creating lists or strings which are larger than `2147483647` will do whatever the golang runtime
// would do. (Which probably is a memory allocation error, and aborting the program.)
//
//	DUMP * "A" 2147483648 #=> might work, depending on the OS
//
// Other types are invalid:
//
//	DUMP * TRUE 34  #!! error: invalid type
func multiply(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	// It just so happens that all three multiply cases need integers as the second argument, so
	// just do the coercion before the typecheck.
	rhs, err := executeToInt(args[1])
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		return Integer(int(lhs) * rhs), nil

	case String:
		if rhs < 0 {
			return nil, fmt.Errorf("negative replication amount for a string in '*': %d", rhs)
		}

		return String(strings.Repeat(string(lhs), rhs)), nil

	case List:
		if rhs < 0 {
			return nil, fmt.Errorf("negative replication amount for a list in '*': %d", rhs)
		}

		return slices.Repeat(lhs, rhs), nil

	default:
		return nil, fmt.Errorf("invalid type given to '*': %T", lhs)
	}
}

// divide divides an integer by another. It returns an error for other types, or if the second
// argument is zero.
//
// ## Examples
//
//	DUMP / 123 "3"       #=> 41
//	DUMP / 12 TRUE       #=> 12
//
// ## Undefined Behaviour
// Overflowing operations on `Integer`s just wraparound. (This is only happens when the minimum
// integer is divided by `-1`.)
//
//	DUMP / (- ~9223372036854775807 1) ~1 #=> -9223372036854775808
//
// Division by zero is an error:
//
//	DUMP / 123 0         #=> error: zero divisor given
//
// Other types are invalid:
//
//	DUMP / TRUE 34  #!! error: invalid type
func divide(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInt(args[1])
		if err != nil {
			return nil, err
		}

		if rhs == 0 {
			return nil, errors.New("zero divisor given to '/'")
		}

		return Integer(int(lhs) / rhs), nil

	default:
		return nil, fmt.Errorf("invalid type given to '/': %T", lhs)
	}
}

// remainder gets the remainder of the first argument and the second. It returns an error for other
// types, or if the second argument is zero.
//
// ## Examples
//
//	DUMP % 34 "12"       #=> 10
//	DUMP % 12 TRUE       #=> 0
//
// ## Undefined Behaviour
// Overflowing operations on `Integer`s just wraparound. (This is only happens when the minimum
// integer is modulo'd by `-1`.)
//
//	DUMP % (- ~9223372036854775807 1) ~1 #=> 0
//
// Modulo by zero is an error:
//
//	DUMP % 123 0         #=> error: zero divisor given
//	DUMP % 123 NULL      #=> error: zero divisor given
//
// Modulo by a negative number is handled by rounding towards zero:
//
//	DUMP % 5 ~2  #=> 1
//	DUMP % ~5 ~2 #=> -1
//
// Other types are invalid:
//
//	DUMP % TRUE 34  #!! error: invalid type
func remainder(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInt(args[1])
		if err != nil {
			return nil, err
		}

		if rhs == 0 {
			return nil, errors.New("zero divisor given to '%'")
		}

		return Integer(int(lhs) % rhs), nil

	default:
		return nil, fmt.Errorf("invalid type given to '%%': %T", lhs)
	}
}

// exponentiate raises the first argument to the power of the second, or joins lists. It returns an
// error for other types, if an integer is raised to a negative power, or if the list contains types
// which cannot be converted to strings (such as `BLOCK`'s return value).
//
// ## Examples
//
//	DUMP ^ 3 "12"       #=> 531441
//	DUMP ^ 12 TRUE      #=> 12
//	DUMP ^ (+@123) TRUE #=> "1true2true3"
//	DUMP ^ @ ":"        #=> ""
//
// ## Undefined Behaviour
// Overflowing operations on `Integer`s are "saturating"---they'll stop at the minimum or maximum
// value for integers.
//
//	DUMP ^ 12  30 #=> 9223372036854775807
//	DUMP ^ ~12 31 #=> -9223372036854775808
//
// Negative integers given to the exponent result in errors:
//
//	DUMP ^ 12 ~3  #!! error: negative exponent
//
// Like other operations, Creating strings which are larger than `2147483647` will do whatever the
// golang runtime would do. (Which probably is a memory allocation error, and aborting the program):
//
//	DUMP ^ (list of 2147483647 elements) ":" #=> might work, depending on the OS
//
// Other types are invalid:
//
//	DUMP ^ TRUE 34  #!! error: invalid type
func exponentiate(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := executeToInt(args[1])
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

		joined, err := lhs.Join(sep) // Join can fail if the list contains Block return value.
		if err != nil {
			return nil, err
		}

		return String(joined), nil

	default:
		return nil, fmt.Errorf("invalid type given to '^': %T", lhs)
	}
}

// compare is a helper method for lessThan and greaterThan. It returns a negative, zero, or positive
// integer depending on whether lhs is less than, equal to, or greater than the second. The
// functionName argument is just used for error messages if an invalid type is provided.
func compare(lhs, rhs Value, functionName rune) (int, error) {
	switch lhs := lhs.(type) {
	case Integer:
		rhs, err := rhs.ToInt()
		if err != nil {
			return 0, err
		}

		// Subtraction actually is all that's needed for integers.
		return int(lhs) - rhs, nil

	case String:
		rhs, err := rhs.ToString()
		if err != nil {
			return 0, err
		}

		// strings.Compare does lexicographical comparisons
		return strings.Compare(string(lhs), rhs), nil

	case Boolean:
		rhs, err := rhs.ToBool()
		if err != nil {
			return 0, err
		}

		// Just manually enumerate all the cases for booleans.
		if !bool(lhs) && rhs {
			return -1, nil
		} else if bool(lhs) && !rhs {
			return 1, nil
		} else {
			return 0, nil
		}

	case List:
		rhs, err := rhs.ToSlice()
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
//
// ## Examples
//
//	DUMP < 10  "2"       #=> false
//	DUMP < "10" 2        #=> true    (converts arg2 to a string, deos lexicographical comparisons)
//	DUMP < FALSE 123     #=> true    (converts arg2 to a boolean)
//	DUMP < (+@123) 13    #=> true    (`2 < 3`)
//	DUMP < (+@123) 12    #=> false   (first argument's length is larger)
//
// ## Undefined Behaviour
// All forms of undefined behaviour in `<` have errors associated with them:
//
//	DUMP < (BLOCK foo) 34   #!! error: invalid type
//	DUMP < ,(BLOCK foo) 34  #!! error: invalid type (even within lists, you cant use `BLOCK`s)
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
//
// See lessThan for examples and undefined behaviour.
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
//
// ## Examples
//
//	DUMP ? 10 10       #=> true
//	DUMP ? 10 "10"     #=> false     (don't coerce)
//	DUMP ? FALSE NULL  #=> false
//	DUMP ? " hi" "hi"  #=> false
//
// ## Undefined Behaviour
// As an extension, `?` called with a `Block` is supported, and returns whether the arguments
// contain the *exact same bodies*
//
//	DUMP ? (BLOCK foo) (BLOCK foo)             #=> true
//	DUMP ? (BLOCK + 1 bar) (BLOCK + 1 bar)     #=> true
//	DUMP ? (BLOCK + 0 + 1 bar) (BLOCK + 1 bar) #=> false, even though semantically the same
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
//
// ## Examples
//
//	DUMP & 4 "hi"        #=> "hi"
//	DUMP & 0  "hi"       #=> 0
//	DUMP & 0 (QUIT 34)   #=> 0         (the other argument isn't even evaluated)
//	DUMP & 4 (QUIT 34)   # (exit status 34)
//	: & 4 (BLOCK foo)    # (works, `&`'s second argument can be a BLOCK.)
//
// ## Undefined Behaviour
// Types which can't be converted to booleans yield an error:
//
//	: & (BLOCK foo) 34   #!! error: cant convert to a boolean
func and(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	isTruthy, err := lhs.ToBool()
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
//
// ## Examples
//
//	DUMP | 4 "hi"        #=> 4
//	DUMP | 0  "hi"       #=> "hi"
//	DUMP | 4 (QUIT 34)   #=> 4         (the other argument isn't even evaluated)
//	DUMP | 0 (QUIT 34)   # (exit status 34)
//	: | 4 (BLOCK foo)    # (works, `|`'s second argument can be a BLOCK.)
//
// ## Undefined Behaviour
// Types which can't be converted to booleans yield an error:
//
//	: | (BLOCK foo) 34   #!! error: cant convert to a boolean
func or(args []Value) (Value, error) {
	lhs, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	isTruthy, err := lhs.ToBool()
	if err != nil {
		return nil, err
	}

	if !isTruthy {
		return args[1].Execute()
	}
	return lhs, nil
}

// then evaluates the first argument, then evaluates and returns the second argument.
//
// ## Examples
//
//	DUMP ; 3 4                     #=> 4
//	DUMP ; (= a 4) a               #=> a
//	; (= a BLOCK + 3 4) (CALL a)   #=> 7
//
// ; (BLOCK foo) (BLOCK bar)      # (ok, both arguments can be a `BLOCK`.)
func then(args []Value) (Value, error) {
	if _, err := args[0].Execute(); err != nil {
		return nil, err
	}

	return args[1].Execute()
}

// assign is used to assign values to variables. The first argument must be a Variable, or an error
// is returned. The second argument is evaluated, and after assignment is returned.
//
// ## Examples
//
//	DUMP = foo 34            #=> 34   (returns itself)
//	; (= foo 34) (DUMP foo)  #=> 34   (assigns for future use)
//	= foo BLOCK bar          # (works, you can assign blocks)
//
// ## Undefined Behaviour
// All forms of undefined behaviour within `=` yield errors:
//
//	= 12 34 #!! error: can only assign variables
func assign(args []Value) (Value, error) {
	// go syntax for "attempt to cast to a Variable pointer". If `args[0]` isn't a variable, then
	// `ok` will be false, which we can check.
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
//
// ## Examples
//
//	; = i = j 0 : WHILE (> 4 = i + i 1) (OUTPUT i)  #=> 1␤2␤3␤
//	DUMP WHILE FALSE 34                             #=> null
//	DUMP WHILE FALSE QUIT 34                        #=> null (doesn't run the body)
//	: WHILE FALSE BLOCK 34                          # (works, `BLOCK` is allowed as the body)
func while(args []Value) (Value, error) {
	// "loop forever" loops in golang are `for { ... }`
	for {
		condition, err := executeToBool(args[0])
		if err != nil {
			return nil, err
		}

		if !condition {
			break
		}

		// Ignore the return value of the body, but return an error if there is one.
		if _, err = args[1].Execute(); err != nil {
			return nil, err
		}
	}

	return Null{}, nil
}

/**************************************************************************************************
 *                                                                                                *
 *                                            Arity 3                                             *
 *                                                                                                *
 **************************************************************************************************/

// if_ evaluates and returns the second argument if the first is truthy; if it's falsey, if_
// evaluates and returns the third argument instead.
//
// ## Examples
//
//	DUMP IF 0 1 2                    #=> 2
//	DUMP IF "yes" 1 2                #=> 1
//	DUMP IF 0 (QUIT 3) 2             #=> 2 (doesn't execute the other branches)
//	IF FALSE (BLOCK foo) (BLOCK bar) # (works, `IF` can accept blocks as 2nd and 3rd args)
//
// ## Undefined Behaviour
// All forms of undefined behaviour within `IF` yield errors:
//
//	IF (BLOCK foo) 3 4 #!! error: cant convert to a boolean
func if_(args []Value) (Value, error) {
	condition, err := executeToBool(args[0])
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
//
// ## Examples
//
//	DUMP GET "" 0 0         # => ""
//	DUMP GET "abcde" 2 2    # => "cd"
//	DUMP GET "abcde" 2 0    # => ""
//	DUMP GET "abcde" 5 0    # => ""
//	DUMP GET "abcde" 4 1    # => "e"
//
//	DUMP GET @ 0 0          # => []
//	DUMP GET (+@12345) 2 2  # => [3, 4]
//	DUMP GET (+@12345) 2 0  # => []
//	DUMP GET (+@12345) 5 0  # => []
//	DUMP GET (+@12345) 4 1  # => [5]
//
// ## Undefined Behaviour
// All forms of undefined behaviour within `GET` yield errors:
//
//	DUMP GET "abcde" 5 1       #!! error, string index out of bounds
//	DUMP GET "abcde" ~1 1      #!! error, string negative start
//	DUMP GET "abcde" 1 ~1      #!! error, string negative length
//
//	DUMP GET (+@"abcde") 5 1   #!! error, list index out of bounds
//	DUMP GET (+@"abcde") ~1 1  #!! error, list negative start
//	DUMP GET (+@"abcde") 1 ~1  #!! error, list negative length
//
//	DUMP GET TRUE 1 2          #!! error, invalid type
func get(args []Value) (Value, error) {
	collection, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	// Get the starting index, returning an error if it's negative
	start, err := executeToInt(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'GET': %d", start)
	}

	// Get the length, returning an error if it's negative
	length, err := executeToInt(args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to GET: '%d'", length)
	}

	// Get the stop index, i.e. where the substring/sublist will end
	stop := start + length

	switch collection := collection.(type) {
	case String:
		if len(collection) < stop {
			return nil, fmt.Errorf("string index out of bounds for 'GET': %d < %d", len(collection), stop)
		}

		return collection[start:stop], nil

	case List:
		if len(collection) < stop {
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
//
// ## Examples
//
//	DUMP SET "" 0 0 "Hello"  # => "Hello"
//	DUMP SET "abcd" 2 1 "!"  # => "ab!d"
//	DUMP SET "abcd" 2 0 "!"  # => "ab!cd"
//	DUMP SET "abcd" 1 2 TRUE # => "atrued"
//	DUMP SET "abcd" 0 2 @    # => "cd"
//
//	DUMP SET @ 0 0 "Hello"        # => ["H", "e", "l", "l", "o"]
//	DUMP SET (+@1234) 2 1 ,9      # => [1, 2, 9, 4]
//	DUMP SET (+@1234) 2 0 "!"     # => [1, 2, "!", 3, 4]
//	DUMP SET (+@1234) 1 2 (+@789) # => [1, 7, 8, 9, 4]
//	DUMP SET (+@1234) 0 2 @       # => [3, 4]
//
// ## Undefined Behaviour
// The following forms of undefined behaviour within `SET` yield errors:
//
//	DUMP SET "abcde" 5 1 "foo"      #!! error, string index out of bounds
//	DUMP SET "abcde" ~1 1 "foo"     #!! error, string negative start
//	DUMP SET "abcde" 1 ~1 "foo"     #!! error, string negative length
//
//	DUMP SET (+@"abcde") 5 1 "foo"  #!! error, list index out of bounds
//	DUMP SET (+@"abcde") ~1 1 "foo" #!! error, list negative start
//	DUMP SET (+@"abcde") 1 ~1 "foo" #!! error, list negative length
//
//	DUMP SET TRUE 1 2          #!! error, invalid type
//
// Creating lists or strings which are larger than `2147483647` will do whatever the golang runtime
// would do. (Which probably is a memory allocation error, and aborting the program.)
//
//	DUMP SET "ABC" 0 1 "<2147483647-character-long string>" #=> might work, depending on the OS
func set(args []Value) (Value, error) {
	collection, err := args[0].Execute()
	if err != nil {
		return nil, err
	}

	// Get the starting index, returning an error if it's negative
	start, err := executeToInt(args[1])
	if err != nil {
		return nil, err
	}
	if start < 0 {
		return nil, fmt.Errorf("negative start given to 'SET': %d", start)
	}

	// Get the length, returning an error if it's negative
	length, err := executeToInt(args[2])
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, fmt.Errorf("negative length given to 'SET': %d", length)
	}

	// Get the stop index, i.e. where the substring/sublist to replace will end
	stop := start + length

	switch collection := collection.(type) {
	case String:
		if len(collection) < stop {
			return nil, fmt.Errorf("string index out of bounds for 'SET': %d < %d", len(collection), stop)
		}

		replacement, err := executeToString(args[3])
		if err != nil {
			return nil, err
		}

		// Use a string builder for efficiency's sake
		var builder strings.Builder
		builder.WriteString(string(collection[:start]))
		builder.WriteString(replacement)
		builder.WriteString(string(collection[stop:]))
		return String(builder.String()), nil

	case List:
		if len(collection) < stop {
			return nil, fmt.Errorf("list index out of bounds for 'SET': %d < %d", len(collection), stop)
		}

		replacement, err := executeToSlice(args[3])
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

// eval converts its argument to a string, and then evaluates that as Knight source code.
//
// ## Examples
//
//	; = foo 34 : DUMP EVAL + "fo" "o"   #=> 34
func eval(args []Value) (Value, error) {
	sourceCode, err := executeToString(args[0])
	if err != nil {
		return nil, err
	}

	return Evaluate(sourceCode)
}

// system converts its argument to a string, and then evaluates that as a shell command, returning
// the stdout of it (less its trailing newline)
//
// ## Examples
//
// DUMP ` "ls" #=> "README.md\ngo\ngo.mod\nknight\nmain.go"
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
	stdout, err := exec.Command(shell, "-c", shellCommand).Output()
	if err != nil {
		return nil, err
	}

	// Delete the last `\n`, `\r`, or `\r\n` to be like `PROMPT`.
	if len(stdout) != 0 && stdout[len(stdout)-1] == '\n' {
		stdout = stdout[:len(stdout)-1]
	}

	if len(stdout) != 0 && stdout[len(stdout)-1] == '\r' {
		stdout = stdout[:len(stdout)-1]
	}

	// Return the stdout
	return String(stdout), nil
}
