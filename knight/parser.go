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
// This parses Knight programs in terms of "rune"s (golang speak for unicode codepoints), instead of
// bytes, as an extension: Our implementation supports all of Unicode, in addition to the restricted
// subset that Knight requires.
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

	// (You have to use `string()` to convert from a `[]rune` to a `string`.)
	return string(p.source[start:p.index])
}

//
// Functions used within ParseNextValue as arguments to TakeWhile.
// 
func isntNewLine(r rune) bool             { return r != '\n' }
func isDigit(r rune) bool                 { return '0' <= r && r <= '9' }
func isVariableStart(r rune) bool         { return unicode.IsLower(r) || r == '_' }
func isVariableBody(r rune) bool     { return isVariableStart(r) || unicode.IsNumber(r) }
func isWordFunctionCharacter(r rune) bool { return unicode.IsUpper(r) || r == '_' }
func isWhitespace(r rune) bool {
	// Note: The parenthesis are also included here because they may also safely be considered
	// whitespace by implementations. (There's an optional extension where implementations can give
	// syntax errors on unbalanced parenthesis, but this implementation isn't doing that.)
	return unicode.IsSpace(r) || r == '(' || r == ')'
}

// ParseNextValue returns the next Value in the source code. EndOfInput is returned if there's no
// Values left. Syntax errors (such as missing an ending quote) are also returned.
func (p *Parser) ParseNextValue() (Value, error) {
	// If we're at the end, return EndOfInput.
	if p.IsAtEnd() {
		return nil, EndOfInput
	}

	// Determine what to do solely based upon the next character---Knight's specs are designed in
	// such a way that you can unambiguously determine what to do solely based on the next character
	// in the input stream.
	c := p.Peek()

	// Whitespace, delete it and parse again.
	if isWhitespace(c) {
		p.Advance()
		return p.ParseNextValue()
	}

	// Comment, delete it and parse again.
	if c == '#' {
		_ = p.TakeWhile(isntNewLine) // ignore the comment line that was parsed
		return p.ParseNextValue()
	}

	// Integers
	if isDigit(c) {
		integer, _ := strconv.Atoi(p.TakeWhile(isDigit))
		return Integer(integer), nil
	}

	// Variables
	if isVariableStart(c) {
		return NewVariable(p.TakeWhile(isVariableBody)), nil
	}

	// Strings
	if c == '\'' || c == '"' {
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
	}

	//
	// Everything else is a function, or invalid (which we check for below).
	//

	// Delete the function name out of the input stream
	if isWordFunctionCharacter(c) {
		_ = p.TakeWhile(isWordFunctionCharacter) // ignore the remainder of the word function
	} else {
		p.Advance()
	}

	startIndex := p.index // used for syntax error messages

	// Get the function definition; If it doesn't exist, then we've been given an invalid token.
	function, ok := KnownFunctions[c]
	if !ok {
		return nil, fmt.Errorf("[line %d] unknown token start: %c", p.linenoAt(startIndex), c)
	}

	// Create a slice with enough room to store all the arguments.
	arguments := make([]Value, function.arity)

	// Parse each argument to the function, returning any errors that might occur. As a special case,
	// we handle EndOfInput errors specially, which allows us to provide a better error message when
	// arguments to a function are missing.
	for i := 0; i < function.arity; i++ {
		var err error

		arguments[i], err = p.ParseNextValue()
		if err != nil {
			// Special case: If the error was EndOfInput, provide a better error message.
			if err == EndOfInput {
				err = fmt.Errorf(
					"[line %d] missing argument %d for function %q",
					p.linenoAt(startIndex),
					i+1,
					function.name,
				)
			}

			return nil, err
		}
	}

	return NewAst(function, arguments), nil
}
