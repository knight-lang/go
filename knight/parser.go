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

// peek returns the next rune in the source without consuming it. Returns `'\0'` at end of file.
func (p *Parser) peek() rune {
	if p.isEOF() {
		return '\000'
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


// EndOfInput indicates that Parser.NextValue was called when the input source was empty.
//
// This is a user error: They either provided a program which was exclusively whitespace/comments,
// or didn't provide enough arguments to a function (eg `DUMP + 1`).
var EndOfInput = errors.New("source was empty")

// NextValue gets the next Value in the input source, returning an error åt end of file or when a
// token has a syntax error.
func (p *Parser) NextValue() (Value, error) {
	// Look at the next rune in the input stream. (If we're at EOF, this is NUL.)
	c := p.peek()

	// Determine what to do based on that character
	switch {
	// End of file, return an EndOfInput error.
	case c == '\000':
		return nil, EndOfInput

	// Whitespace, delete it and try again.
	//
	// The `:` function is considered whitespace as an optimization, as all it does is execute its sole
	// argument and return it, functionally identical to if it didn't exist.
	//
	// The parenthesis are deleted, as the Knight specification says that all conforming programs must
	// balance any parenthesis they have around whole expressions (eg `(+ 1 2)` is valid but `+ (1 2)`
	// isn't). This allows implementations that want to provide syntax errors to warn users that their
	// parenthesis aren't. However, this implementation doesn't do that optional extension, and instead
	// just pretends like they're whitespace too.
	case unicode.IsSpace(c), c == ':' || c =='(' || c == ')':
		p.advance()
		return p.NextValue()

	// Comment. Delete it, and then try parsing again.
	case c == '#':
		_ = p.takeWhile(func (r rune) bool {
			return r != '\n'
		})
		return p.NextValue()

	// Integers
	case '0' <= c && c <= '9':
		integerString := p.takeWhile(func(r rune) bool { return '0' <= r && r <= '9' })
		integer, _ := strconv.Atoi(integerString)
		return Integer(integer), nil

	// Variables
	case unicode.IsLower(c), c == '_':
		variableName := p.takeWhile(func(r rune) bool {
			return r == '_' || unicode.IsLower(r) || unicode.IsDigit(r)
		})
		return NewVariable(variableName), nil

	// Strings
	case c == '\'' || c == '"':
		startIndex := p.index // for error msgs
		p.advance() // Consume the starting quote.

		// Read until we hit the ending quote, but don't actually consume it.
		contents := p.takeWhile(func(r rune) bool { return r != c })

		// If we reached end of file, that means we never found the ending quote.
		if p.isEOF() {
			return nil, fmt.Errorf("[line %d] unterminated %q string", p.linenoAt(startIndex), c)
		}

		// Consume the ending quote, and return the contents of the string.
		p.advance()
		return String(contents), nil

	// TRUE and FALSE
	case c == 'T' || c == 'F': 
		p.takeWhile(isUpper)
		return Boolean(c == 'T'), nil

	// NULL
	case c == 'N':
		p.takeWhile(isUpper)
		return &Null{}, nil

	// @, ie empty list
	case c == '@':
		p.advance()
		return &List{}, nil

	// Everything else is a function (or a parse error)
	default:
		if isUpper(c) {
			p.takeWhile(isUpper)
		} else {
			p.advance()
		}

		startIndex := p.index // for error msgs

		// Get the function definition. If it doesn't exist, then the user's given us an invalid program,
		// and we error out.
		function, ok := KnownFunctions[c]
		if !ok {
			return nil, fmt.Errorf("[line %d] unknown token start: %c", p.linenoAt(startIndex), c)
		}

		arguments := make([]Value, function.arity) // Pre-allocate enough room to store all args.

		// Parse each argument and add them to the `arguments`.
		for i := 0; i < function.arity; i++ {
			// Try to parse the argument
			argument, err := p.NextValue()

			if err != nil {
				if err == EndOfInput {
					// Special case: If the error was `EndOfInput`, provide a better error message.
					return nil, fmt.Errorf("[line %d] missing argument %d for function %q", i + 1, function.name)
				}

				return nil, err
			}

			arguments[i] = argument
		}

		return NewAst(function, arguments), nil
	}
}
