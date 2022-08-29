package knight

import (
	"fmt"
	"strconv"
	"unicode"
)

/// Parser is used to parse `Value`s from source code.
type Parser struct {
	source []rune
	index  int
}

/// NewParser creates a Parser for the given source string
func NewParser(source string) Parser {
	return Parser{source: []rune(source)}
}

func (p *Parser) isEOF() bool {
	return len(p.source) <= p.index
}

// This is a function because we only need to look at itonce
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
		panic("peeking empty parser")
	}

	return p.source[p.index]
}

func (p *Parser) advance() {
	p.index++
}

func (p *Parser) takeWhile(cond func(rune) bool) string {
	start := p.source

	for !p.isEOF() && cond(p.peek()) {
		p.advance()
	}

	return string(start[:len(start)-len(p.source)])
}

func (p *Parser) strip() {
	not_newline := func(r rune) bool { return r != '\n' }
	is_whitespace := func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')' || r == ':'
	}

	for !p.isEOF() {
		p.takeWhile(is_whitespace)

		if p.isEOF() || p.peek() != '#' {
			break
		}

		p.takeWhile(not_newline)
	}
}

func (p *Parser) Parse(e *Environment) (Value, error) {
	is_digit := func(r rune) bool { return '0' <= r && r <= '9' }
	is_lower := func(r rune) bool { return unicode.IsLower(r) || r == '_' }
	is_upper := func(r rune) bool { return unicode.IsUpper(r) || r == '_' }

	// Remove whitespace, and return `nil, nil` if at EOF.
	p.strip()
	if p.isEOF() {
		return nil, nil // nothing's wrong, there's just nothing to parse.
	}

	head := p.peek()

	// Parse numbers.
	if is_digit(head) {
		num, _ := strconv.Atoi(p.takeWhile(is_digit))
		return Number(num), nil
	}

	// Parse identifiers.
	if is_lower(head) {
		name := p.takeWhile(func(r rune) bool { return is_lower(r) || is_digit(r) })

		return e.Lookup(name), nil
	}

	// Parse strings.
	if head == '\'' || head == '"' {
		p.advance() // gobble up the quote
		start_index := p.index
		body := p.takeWhile(func(r rune) bool { return r != head })

		if p.isEOF() {
			return nil, fmt.Errorf("[line %d] unterminated %q string", p.linenoAt(start_index), head)
		}

		p.advance()
		return Text(body), nil
	}

	// Everything else follows the function format, so remove it accordingly.
	if is_upper(head) {
		p.takeWhile(is_upper)
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

	// Fetch the function, if it doesn't exist it means it's an invalid identifier.
	fun := GetFunction(head)
	if fun == nil {
		return nil, fmt.Errorf("[line %d] unknown token start: %q", head)
	}

	// Start index is for error messages.
	start_index := p.index

	// Create the ast, and parse out all its arguments.
	ast := &Ast{
		fun:  fun,
		args: make([]Value, fun.arity),
	}

	for i := 0; i < fun.arity; i++ {
		// `p.Parse` will return one of:
		// - `nil, <err>` if an error occurred
		// - `nil, nil` if nothing could be parsed (i.e. EOF)
		// - `<val>, nil` if no errors occured and it could parse something.
		arg, err := p.Parse(e)

		if err != nil {
			return nil, err
		}

		if arg == nil {
			return nil, fmt.Errorf("[line %d] missing argument %d for function %q",
				p.linenoAt(start_index), i, fun.name)
		}

		ast.args[i] = arg
	}

	return ast, nil
}
