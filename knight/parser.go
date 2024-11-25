package knight

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Parser is used to construct `Value`s from source code.
//
// Note that we parse Knight programs in terms of `rune`s (golang speak for utf-8 characters), not
// bytes. This is technically unnecessary as Knight only requires implementations to support a
// specific subset of ASCII, but we support it to allow our users to use UTF-8.
//
// The strategy we use to parse out knight programs is to examine each character at a time from the
// `source`, which lets us know what to do next. Knight's spec is designed so that the next
// expression/value/token is always unambiguously determined by the first non-whitespace non-comment
// character. (e.g. if we read a `D`, we know we're going to be executing `DUMP`.)
type Parser struct {
	// source is the contents of the program.
	source []rune

	// index is the next `rune` to look at.
	index int
}

// NewParser creates a `Parser` for the given source string.
func NewParser(source string) Parser {
	return Parser{
		source: []rune(source),
		index:  0,
	}
}

// isEOF returns whether or not the parser is at the end of the file.
func (p *Parser) isEOF() bool {
	return len(p.source) <= p.index
}

// linenoAt returns the line numeber at the given index.
//
// This is a function (and not an field on `Parser` like one normally might do) because we only need
// to determine the line number when an error is happening, at which point efficiency is irrelevant.
func (p *Parser) linenoAt(index int) int {
	lineno := 0

	for i := 0; i < index; i++ {
		if p.source[i] == '\n' {
			lineno++
		}
	}

	return lineno
}

// peek returns the next `rune` in the input stream without consuming it.
//
// It'll panic if the input stream is at the end (in which case there's nothing left to read).
func (p *Parser) peek() rune {
	if p.isEOF() {
		bug("peeking past the end of the parser")
	}

	return p.source[p.index]
}

// advance consumes the next character.
//
// It'll panic if the input stream is at the end (in which case there's nothing left to read).
func (p *Parser) advance() {
	if p.isEOF() {
		bug("advancing past the end of the parser")
	}

	p.index++
}

// takeWhile returns a string containing all the runes from the front of `source` for which `cond`
// returned true. If `cond` is never true, the empty string is returned.
func (p *Parser) takeWhile(cond func(rune) bool) string {
	start := p.index

	for !p.isEOF() && cond(p.peek()) {
		p.advance()
	}

	return string(p.source[start:p.index])
}

// stripWhitespaceAndComments removes leading whitespace and comments from `source`, ensuring after
// it returns that the stream is either at EOF (which indicates a syntax error on the user's part),
// or the next rune is a non-whitespace non-comment character.
func (p *Parser) stripWhitespaceAndComments() {
	// Functions to be used with `takeWhile`
	//
	// They could be their own `func xxx(...)` instead of local variables, but this way they're local
	// to `stripWhitespaceAndComments`.
	isntNewline := func(r rune) bool {
		return r != '\n'
	}
	isWhitespace := func(r rune) bool {
		// The Knight spec technically only requires spaces, tabs, newlines, and carriage returns to
		// be recognized as whitespace. However, since all other whitespace characters aren't a part
		// of the Knight encoding (i.e. valid Knight programs cannot contain them anywhere), we can
		// strip _all_ space characters an extension.
		//
		// The `:` function is also parsed as whitespace, as all it does is execute its argument and
		// return it, which is functionally identical to if it didn't exist. So, we ignore it. (`:`
		// really only exists to make the last line in long chains of `;` be visually aligned, to make
		// Knight programs much easier to write.)
		//
		// The braces might look weird. The Knight spec requires conforming programs to balance their
		// `(` and `)` around whole expressions (eg `(+ 1 2)` is valid but `+ (1 2)` is not); Outside
		// of that, they have no meaning. Since all valid Knight programs will have balanced `(` and
		// `)`, we can completely ignore them. (We'd pay more attention to them if we want to provide
		// error messages to users that their programs are ill-formed, but since we don't need to
		// worry about ill-formed programs (spec-compliant programs can can assume all inputs programs
		// are valid), we ignore them.)
		return unicode.IsSpace(r) || r == ':' || r == '(' || r == ')'
	}

	for {
		// First, delete all leading whitespace.
		p.takeWhile(isWhitespace)

		// Now, if the next character is a `#`, then delete the comment at start again.
		if !p.isEOF() && p.peek() == '#' {
			p.takeWhile(isntNewline)
			continue
		}

		// The next character wasn't a `#`, which means we deleted all whitespace and comments. exit.
		return
	}
}

// EndOfInput indicates that `Parser.Parse` was called when the input stream was empty.
//
// This is a user error: They either provided a program which was exclusively whitespace/comments,
// or didn't provide enough arguments to a function (eg `DUMP + 1`).
var EndOfInput = errors.New("source was empty")

// Parse returns the next `Value` for the parser.
//
// If there's nothing left to parse, `EndOfInput` is returned.
func (p *Parser) Parse() (Value, error) {
	// Functions to be used with `takeWhile`
	//
	// They could be their own `func xxx(...)` instead of local variables, but this way they're local
	// to `Parse`.
	isDigit := func(r rune) bool { return '0' <= r && r <= '9' }
	isLower := func(r rune) bool { return unicode.IsLower(r) || r == '_' }
	isUpper := func(r rune) bool { return unicode.IsUpper(r) || r == '_' }
	isLowerOrDigit := func(r rune) bool { return isLower(r) || isDigit(r) }

	// Remove leading whitespace and comments
	p.stripWhitespaceAndComments()

	// If there's nothing left to parse, return an error.
	if p.isEOF() {
		return nil, EndOfInput
	}

	// Now, try to parse the next token based on the first charcater.
	next := p.peek()
	switch {
	// Parse integers.
	case isDigit(next):
		num, _ := strconv.Atoi(p.takeWhile(isDigit))
		return Integer(num), nil

	// Parse variables.
	case isLower(next):
		return NewVariable(p.takeWhile(isLowerOrDigit)), nil

	// Parse strings.
	case next == '\'' || next == '"':
		p.advance() // Consume the starting quote.

		// Keep the start index for the error message
		startIndex := p.index

		// Read until we hit the ending quote, but don't actually consume it.
		contents := p.takeWhile(func(r rune) bool {
			return r != next
		})

		// If we reached end of file, that means we never found the ending quote.
		if p.isEOF() {
			return nil, fmt.Errorf("[line %d] unterminated %q string", p.linenoAt(startIndex), next)
		}

		// Consume the ending quote, and return the contents of the string.
		p.advance()
		return String(contents), nil

	// Last up is functions. Here we strip out the function name, and then exit the switch statement
	// so we can parse the arguments to the function. (We check for invalid function names below.)
	case isUpper(next):
		p.takeWhile(isUpper)
	default:
		p.advance()
	}

	// Special-case "function literals": Functions which take no arguments and always return the same
	// value (`TRUE`, `FALSE`, `NULL`, and the empty array `@`). They can be parsed as literals as an
	// optimization, and not have to go through the whole "AST" shebang below.
	switch next {
	case 'T':
		return Boolean(true), nil
	case 'F':
		return Boolean(false), nil
	case 'N':
		return &Null{}, nil
	case '@':
		return &List{}, nil
	}

	// Keep the start index for the error message
	startIndex := p.index

	// Get the function definition. If it doesn't exist, then the user's given us an invalid program,
	// and we error out.
	function, ok := KnownFunctions[next]
	if !ok {
		return nil, fmt.Errorf("[line %d] unknown token start: %q", p.linenoAt(startIndex), next)
	}

	ast := &Ast{
		function:  function,
		arguments: make([]Value, function.arity), // Pre-allocate enough room to store all args.
	}

	// Parse each argument and add them to the `arguments`.
	for i := 0; i < ast.function.arity; i++ {
		// Try to parse the argument
		argument, err := p.Parse()

		// If there were no problems parsing, then just assign and keep going
		if err == nil {
			ast.arguments[i] = argument
			continue
		}

		// Uh oh! There was a problem parsing the argument.

		// Special case: If the error was `EndOfInput`, provide a better error message.
		if err == EndOfInput {
			err = fmt.Errorf(
				"[line %d] missing argument %d for function %q",
				p.linenoAt(startIndex),
				i,
				ast.function.name,
			)
		}

		// Return the error
		return nil, err
	}

	// Create the AST and return it.
	return ast, nil
}
