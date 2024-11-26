package knight

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// EndOfInput indicates that Parser.ParseNextValue was called when the input source was empty.
//
// This is a user error: They either provided a program which was exclusively whitespace/comments,
// or didn't provide enough arguments to a function (eg `DUMP + 1`).
var EndOfInput = errors.New("source was empty")

// Parser is used to construct Values from source code.
//
// This parses Knight programs in terms of `rune`s (golang speak for utf-8 characters), not bytes.
// This is technically unnecessary as Knight only requires implementations to support a specific
// subset of ASCII, but we're supporting them as an extension.
//
// The strategy we use to parse out knight programs is to examine each rune at a time from the
// source, which lets us know what to do next. Knight's spec is designed so that the next expression
// is always unambiguously determined by the first non-whitespace non-comment rune. (e.g. if
// we read a `D`, we know we're going to be executing `DUMP`.)
type Parser struct {
	source []rune // the contents of the program.
	index  int    // index of the next rune to look at.
}

// NewParser creates a Parser for the given source string.
func NewParser(source string) Parser {
	return Parser{
		source: []rune(source),
		index:  0,
	}
}

// IsAtEnd returns whether the parser is at the end of its stream.
func (p *Parser) IsAtEnd() bool {
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

// Peek returns the next rune without consuming it. It panics at the end of the source.
func (p *Parser) Peek() rune {
	if p.IsAtEnd() {
		panic("peeked past end of source")
	}

	return p.source[p.index]
}

// Advance consumes the next rune. It panics at the end of the source
func (p *Parser) Advance() {
	if p.IsAtEnd() {
		panic("advanced past end of source")
	}

	p.index++
}

// TakeWhile consumes runes from the source while the condition is true, and then returns them. If
// the condition is never true, an empty string is returned.
func (p *Parser) TakeWhile(condition func(rune) bool) string {
	start := p.index

	for !p.IsAtEnd() && condition(p.Peek()) {
		p.Advance()
	}

	// (You have to use string() to convert from []rune to a string.)
	return string(p.source[start:p.index])
}

/**
 * Functions used within `TakeWhile`
 **/
func isntNewLine(r rune) bool             { return r != '\n' }
func isDigit(r rune) bool                 { return '0' <= r && r <= '9' }
func isVariableStart(r rune) bool         { return unicode.IsLower(r) || r == '_' }
func isVariableCharacter(r rune) bool     { return isVariableStart(r) || unicode.IsNumber(r) }
func isWordFunctionCharacter(r rune) bool { return unicode.IsUpper(r) || r == '_' }
func isWhitespace(r rune) bool {
	// Note: The parenthesis are also included here because they may also safely be considered
	// whitespace by implementations. (There's an optional extension where implementations can give
	// syntax errors on unbalanced parenthesis if they want, but this implementation isn't doing that.)
	return unicode.IsSpace(r) || r == '(' || r == ')'
}

// ParseNextValue returns the next Value in the source code. EndOfInput is returned if there's no
// Values left. Syntax errors (such as missing an ending quote) are also returned.
func (p *Parser) ParseNextValue() (Value, error) {
	// If we're at the end, then just return.
	if p.IsAtEnd() {
		return nil, EndOfInput
	}

	// Determine what to do solely based upon the next character. (Knight's designed in such a way
	// that you can always unambiguously parse based on just the next character in the input stream.)
	c := p.Peek()
	switch {
	// Whitespace, delete it and go again.
	case isWhitespace(c):
		p.Advance()
		return p.ParseNextValue()

	// Comment, delete it and go again.
	case c == '#':
		_ = p.TakeWhile(isntNewLine)
		return p.ParseNextValue()

	// Integers
	case isDigit(c):
		integer, _ := strconv.Atoi(p.TakeWhile(isDigit))
		return Integer(integer), nil

	// Variables
	case isVariableStart(c):
		return NewVariable(p.TakeWhile(isVariableCharacter)), nil

	// Strings
	case c == '\'' || c == '"':
		startIndex := p.index // for error msgs
		p.Advance()           // Consume the starting quote.

		// Read until we hit the ending quote, but don't actually consume it.
		contents := p.TakeWhile(func(r rune) bool { return r != c })

		// If we reached end of file, that means we never found the ending quote.
		if p.IsAtEnd() {
			return nil, fmt.Errorf("[line %d] unterminated %q string", p.linenoAt(startIndex), c)
		}

		// Consume the ending quote, and return the contents of the string.
		p.Advance()
		return String(contents), nil

	// Everything else is a function, or invalid (which we check for below).
	default:
		startIndex := p.index // used for error msgs

		// Delete the function name out
		if isWordFunctionCharacter(c) {
			_ = p.TakeWhile(isWordFunctionCharacter)
		} else {
			p.Advance()
		}

		// Get the function definition; If it doesn't exist, then we've been given an invalid token.
		function, ok := KnownFunctions[c]
		if !ok {
			return nil, fmt.Errorf("[line %d] unknown token start: %c", p.linenoAt(startIndex), c)
		}

		// Pre-allocate enough room to store all args.
		arguments := make([]Value, function.arity)

		// Parse each argument and add them to the `arguments`.
		for i := 0; i < function.arity; i++ {
			argument, err := p.ParseNextValue()
			if err != nil {
				// Special case: If the error was `EndOfInput`, provide a better error message.
				if err == EndOfInput {
					err = fmt.Errorf("[line %d] missing argument %d for function %q",
						p.linenoAt(startIndex), i+1, function.name)
				}
				return nil, err
			}

			arguments[i] = argument
		}

		return NewAst(function, arguments), nil
	}
}
