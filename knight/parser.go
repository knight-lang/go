package knight

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Parser is used to construct Values from source code.
//
// This parses Knight programs in terms of `rune`s (golang speak for utf-8 characters), not bytes.
// This is technically unnecessary as Knight only requires implementations to support a specific
// subset of ASCII, but we're supporting them as an extension.
//
// The strategy we use to parse out knight programs is to examine each character at a time from the
// source, which lets us know what to do next. Knight's spec is designed so that the next expression
// is always unambiguously determined by the first non-whitespace non-comment character. (e.g. if
// we read a `D`, we know we're going to be executing `DUMP`.)
type Parser struct {
	source []rune // the contents of the program.
	index  int    // index of the next rune to look at.
}

// NewParser creates a Parser for the given source string.
func NewParser(source string) Parser {
	return Parser{
		source: []rune(source),
		index:  0, // Technically redundant since 0 is the default, but it never hurts to be explicit.
	}
}

// isEOF returns whether the parser is at the end of the file.
func (p *Parser) isEOF() bool {
	return len(p.source) <= p.index
}

// linenoAt returns the line number at the given index. It's used in syntax error messages.
func (p *Parser) linenoAt(index int) int {
	lineno := 0

	for i := 0; i < index; i++ {
		if p.source[i] == '\n' {
			lineno++
		}
	}

	return lineno
}

// peek returns the next rune in the source without consuming it. It panics when at EOF.
func (p *Parser) peek() rune {
	if p.isEOF() {
		panic("peeked past the end of the parser")
	}

	return p.source[p.index]
}

// advance consumes the next character. It panics at EOF.
func (p *Parser) advance() {
	if p.isEOF() {
		panic("advanced past the end of the parser")
	}

	p.index++
}

// takeWhile consumes runes from the source while the condition is true, and then returns them. If
// the condition is never true, an empty string is returned.
func (p *Parser) takeWhile(condition func(rune) bool) string {
	start := p.index

	for !p.isEOF() && condition(p.peek()) {
		p.advance()
	}

	return string(p.source[start:p.index])
}


// Helper functions for `stripWhitespaceAndComments` and `NextValue`
func isWhitespace(r rune) bool { return unicode.IsSpace(r) || r == ':' || r == '(' || r == ')' }
func isntNewline(r rune) bool { return r != '\n' }
func isDigit(r rune) bool { return '0' <= r && r <= '9' }
func isLower(r rune) bool { return unicode.IsLower(r) || r == '_' }
func isUpper(r rune) bool { return unicode.IsUpper(r) || r == '_' }
func isLowerOrDigit(r rune) bool { return isLower(r) || isDigit(r) }


// stripWhitespaceAndComments removes all leading whitespace and comments from source. Colons (`:`)
// and parenthesis (`(` and `)`) are also considered whitespace and are stripped.
//
// The `:` function is considered whitespace as an optimization, as all it does is execute its sole
// argument and return it, functionally identical to if it didn't exist.
//
// The parenthesis are stripped, as the Knight specification says that all conforming programs must
// balance any parenthesis they have around whole expressions (eg `(+ 1 2)` is valid but `+ (1 2)`
// isn't). This allows implementations that want to provide syntax errors to warn users that their
// parenthesis aren't. However, this implementation doesn't do that optional extension, and instead
// just pretends like they're whitespace too.
func (p *Parser) stripWhitespaceAndComments() {
	for {
		// First, delete all leading whitespace.
		p.takeWhile(isWhitespace)

		// Since we've deleted all the whitespace, if the next character isn't a `#` then we're done.
		if p.isEOF() || p.peek() != '#' {
			return
		}

		// Strip out the comment and go again. 
		p.takeWhile(isntNewline)
	}
}

// parseError is a helper function which adds a line number to the start of an error message.
func (p *Parser) parseError(startIndex int, fmt string, rest ...any) error {
	return fmt.Errorf("[line %d] " + fmt, p.linenoAt(startIndex), rest...)
}

// EndOfInput indicates that Parser.NextValue was called when the input source was empty.
//
// This is a user error: They either provided a program which was exclusively whitespace/comments,
// or didn't provide enough arguments to a function (eg `DUMP + 1`).
var EndOfInput = errors.New("source was empty")

// NextValue gets the next Value in the input source, returning an error åt end of file or when a
// token has a syntax error.
func (p *Parser) NextValue() (Value, error) {
	// Remove any leading whitespace and comments
	p.stripWhitespaceAndComments()

	// Get the start index; it's used for error messages so we know what line the parsing started on.
	startIndex := p.index

	// If there's nothing left, then that means we're at the end of the input.
	if p.isEOF() {
		return nil, EndOfInput
	}

	// Now, try to parse the next token based on the first character.
	tokenStart := p.peek()

	// Integers
	if isDigit(tokenStart) {
		num, _ := strconv.Atoi(p.takeWhile(isDigit))
		return Integer(num), nil
	}

	// Variables
	if isLower(tokenStart) {
		return NewVariable(p.takeWhile(isLowerOrDigit)), nil
	}

	// Strings
	if tokenStart == '\'' || tokenStart == '"' {
		p.advance() // Consume the starting quote.

		// Read until we hit the ending quote, but don't actually consume it.
		contents := p.takeWhile(func(r rune) bool {
			return r != tokenStart
		})

		// If we reached end of file, that means we never found the ending quote.
		if p.isEOF() {
			return nil, p.parseError(startIndex," unterminated %q string", tokenStart)
		}

		// Consume the ending quote, and return the contents of the string.
		p.advance()
		return String(contents), nil
	}

	// Last up is functions. Here we strip out the function name, and then exit the switch statement
	// so we can parse the arguments to the function. (We check for invalid function names below.)
	if isUpper(tokenStart) {
		p.takeWhile(isUpper)
	} else {
		p.advance()
	}

	// Special-case "function literals": Functions which take no arguments and always return the same
	// value (`TRUE`, `FALSE`, `NULL`, and the empty array `@`). They can be parsed as literals as an
	// optimization, and not have to go through the whole "AST" shebang below.
	switch tokenStart {
	case 'T':
		return Boolean(true), nil
	case 'F':
		return Boolean(false), nil
	case 'N':
		return &Null{}, nil
	case '@':
		return &List{}, nil
	}

	// Get the function definition. If it doesn't exist, then the user's given us an invalid program,
	// and we error out.
	function, ok := KnownFunctions[tokenStart]
	if !ok {
		return nil, p.parseError(startIndex, "unknown token start: %q", tokenStart)
	}

	arguments := make([]Value, function.arity) // Pre-allocate enough room to store all args.

	// Parse each argument and add them to the `arguments`.
	for i := 0; i < function.arity; i++ {
		// Try to parse the argument
		argument, err := p.NextValue()

		// If there were no problems parsing, then just assign and keep going
		if err == nil {
			arguments[i] = argument
			continue
		}

		// Uh oh! There was a problem parsing the argument.

		// Special case: If the error was `EndOfInput`, provide a better error message.
		if err == EndOfInput {
			err = p.parseError(startIndex, "missing argument %d for function %q", i + 1, function.name)
		}

		// Return the error
		return nil, err
	}

	return NewAst(function, arguments), nil
}
