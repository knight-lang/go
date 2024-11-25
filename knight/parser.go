package knight

import (
	"fmt"
	"strconv"
	"unicode"
)

// Parser is used to construct `Value`s from source code.
type Parser struct {
	source []rune
	index  int
}

// NewParser creates a `Parser` for the given source string.
func NewParser(source string) Parser {
	// TODO: converting the entire source into a list of runes probably isn't ideal.
	return Parser{source: []rune(source)}
}

func (p *Parser) isEOF() bool {
	return len(p.source) <= p.index
}

// This is a function, and not an instance variable, because we only need to determine it
// when an error is happening (at which point, efficiency is irrelevant)
func (p *Parser) linenoAt(index int) int {
	lineno := 0

	for i := 0; i < index; i++ {
		if p.source[i] == '\n' {
			lineno++
		}
	}

	return lineno
}

func (p *Parser) peek() rune {
	if p.isEOF() {
		panic("peeking at an empty parser")
	}

	return p.source[p.index]
}

func (p *Parser) advance() {
	p.index++
}

func (p *Parser) takeWhile(cond func(rune) bool) string {
	start := p.index

	for !p.isEOF() && cond(p.peek()) {
		p.advance()
	}

	return string(p.source[start:p.index])
}

func (p *Parser) strip() {
	isntNewline := func(r rune) bool { return r != '\n' }
	isWhitespace := func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')' || r == ':'
	}

	for {
		// first, strip all leading whitespace.
		p.takeWhile(isWhitespace)

		// Then, if the next character isn't the start of a comment, we're done stripping.
		if p.isEOF() || p.peek() != '#' {
			return
		}

		p.takeWhile(isntNewline)
	}
}

// NothingToParse indicates that `Parser.Parse` was called when no more tokens remained.
var NothingToParse = fmt.Errorf("nothing to parse")

// Parse returns the next Value within the parser's source code.
//
// If there's nothing left to parse, `NothingToParse` is returned.
func (p *Parser) Parse(e *Environment) (Value, error) {
	isDigit := func(r rune) bool { return '0' <= r && r <= '9' }
	isLower := func(r rune) bool { return unicode.IsLower(r) || r == '_' }
	isUpper := func(r rune) bool { return unicode.IsUpper(r) || r == '_' }
	isLowerOrDigit := func(r rune) bool { return isLower(r) || isDigit(r) }

	// Remove whitespace, and return `nil, nil` if at EOF.
	p.strip()
	if p.isEOF() {
		return nil, NothingToParse
	}

	head := p.peek()

	// Parse numbers.
	if isDigit(head) {
		num, _ := strconv.Atoi(p.takeWhile(isDigit))
		return Integer(num), nil
	}

	// Parse identifiers.
	if isLower(head) {
		return e.Lookup(p.takeWhile(isLowerOrDigit)), nil
	}

	// Parse strings.
	if head == '\'' || head == '"' {
		p.advance() // gobble up the quote

		startIndex := p.index
		contents := p.takeWhile(func(r rune) bool { return r != head })

		if p.isEOF() {
			return nil, fmt.Errorf("[line %d] unterminated %q string", p.linenoAt(startIndex), head)
		}

		p.advance()
		return Text(contents), nil
	}

	// Everything else follows the function format, so remove it accordingly.
	if isUpper(head) {
		p.takeWhile(isUpper)
	} else {
		p.advance()
	}

	// Parse literals
	switch head {
	case 'T':
		return Boolean(true), nil
	case 'F':
		return Boolean(false), nil
	case 'N':
		return &Null{}, nil
	case '@':
		return &List{}, nil
	}

	// Start index is for error messages.
	startIndex := p.index

	// Fetch the function, if it doesn't exist it means it's an invalid identifier.
	fun := e.GetFunction(head)
	if fun == nil {
		return nil, fmt.Errorf("[line %d] unknown token start: %q", p.linenoAt(startIndex), head)
	}

	// Create the ast, and parse out all its arguments.
	args := make([]Value, fun.arity)

	// Parse each argument and add them to the `args`.
	for i := 0; i < fun.arity; i++ {
		arg, err := p.Parse(e)

		if err == NothingToParse {
			return nil, fmt.Errorf("[line %d] missing argument %d for function %q",
				p.linenoAt(startIndex), i, fun.name)
		}

		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return NewAst(fun, args), nil
}
