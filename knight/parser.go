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
	// TODO: Converting the entire source into a list of runes probably isn't ideal.
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
	isnt_newline := func(r rune) bool { return r != '\n' }
	is_whitespace := func(r rune) bool {
		return unicode.IsSpace(r) || r == '(' || r == ')' || r == ':'
	}

	for {
		// first, strip all leading whitespace.
		p.takeWhile(is_whitespace)

		// Then, if the next character isn't the start of a comment, we're done stripping.
		if p.isEOF() || p.peek() != '#' {
			return
		}

		p.takeWhile(isnt_newline)
	}
}

// Parse returns the next Value within the parser's source code.
//
// If there's nothing left to parse, `nil` is returned for the `Value`.
func (p *Parser) Parse(e *Environment) (Value, error) {
	is_digit := func(r rune) bool { return '0' <= r && r <= '9' }
	is_lower := func(r rune) bool { return unicode.IsLower(r) || r == '_' }
	is_upper := func(r rune) bool { return unicode.IsUpper(r) || r == '_' }
	is_ident_body := func(r rune) bool { return is_lower(r) || is_digit(r) }

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
		return e.Lookup(p.takeWhile(is_ident_body)), nil
	}

	// Parse strings.
	if head == '\'' || head == '"' {
		p.advance() // gobble up the quote

		start_index := p.index
		contents := p.takeWhile(func(r rune) bool { return r != head })

		if p.isEOF() {
			return nil, fmt.Errorf("[line %d] unterminated %q string", p.linenoAt(start_index), head)
		}

		p.advance()
		return Text(contents), nil
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

	// Parse each argument and add them to the `ast`'s body.
	for i := 0; i < fun.arity; i++ {
		arg, err := p.Parse(e)

		if err != nil {
			return nil, err
		}

		// `arg` is nil when nothing could be parsed. This means an argument was missing.
		if arg == nil {
			return nil, fmt.Errorf("[line %d] missing argument %d for function %q",
				p.linenoAt(start_index), i, fun.name)
		}

		ast.args[i] = arg
	}

	return ast, nil
}
