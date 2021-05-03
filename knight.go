package knight

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

func Run(input string) (Value, error) {
	val, err := Parse(strings.NewReader(input))

	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, errors.New("Nothing to parse")
	}

	return val.Run()
}

func takeWhile(reader *strings.Reader, cond func(rune) bool) (string, error) {
	str := ""

	for {
		r, _, err := reader.ReadRune()

		switch {
		case err == io.EOF:
			return str, nil
		case err != nil:
			return "", err
		case cond(r):
			str += string(r)
		default:
			if err := reader.UnreadRune(); err != nil {
				return "", err
			}

			return str, nil
		}
	}
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r) || strings.ContainsRune("(){}[]:", r)
}

func isAsciiDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isLower(r rune) bool {
	return unicode.IsLower(r) || r == '_'
}

func Parse(reader *strings.Reader) (Value, error) {
	var val Value
	var took string

	r, _, err := reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, err
	}

	switch {
	case r == '#':
		_, err := takeWhile(reader, func(r rune) bool { return r != '\n' })

		if err != nil {
			break
		}

		fallthrough
	case isWhitespace(r):
		return Parse(reader)

	case isAsciiDigit(r):
		took, err = takeWhile(reader, isAsciiDigit)

		if err == nil {
			num, _ := strconv.Atoi(string(r) + took)
			val = Number(num)
		}

	case r == '\'', r == '"':
		took, err = takeWhile(reader, func(r2 rune) bool { return r2 != r })

		if err == nil {
			_, _, err2 := reader.ReadRune()

			if err2 == io.EOF {
				err = errors.New("unterminated string encountered")
			} else {
				val = Text(took)
			}
		}

	case isLower(r):
		took, err = takeWhile(reader, func(r rune) bool { return isLower(r) || unicode.IsDigit(r) })

		if err == nil {
			val = NewVariable(string(r) + took)
		}

	case r == 'T', r == 'F':
		_, err = takeWhile(reader, unicode.IsUpper)

		if err == nil {
			val = Boolean(r == 'T')
		}

	case r == 'N':
		_, err = takeWhile(reader, unicode.IsUpper)

		if err == nil {
			val = &Null{}
		}

	default:
		fun := GetFunction(r)

		if fun == nil {
			err = errors.New("unknown token start: " + string(r))
			break
		}

		if unicode.IsUpper(r) {
			_, err = takeWhile(reader, unicode.IsUpper)

			if err != nil {
				break
			}
		}

		ast := &Ast{fun: fun, args: make([]Value, fun.arity)}
		val = ast

		for i := 0; i < len(ast.args); i++ {
			arg, err := Parse(reader)

			switch {
			case err != nil:
				return nil, err
			case arg == nil:
				return nil, fmt.Errorf("missing argument '%d' for function '%c'", i, fun)
			default:
				ast.args[i] = arg
			}
		}
	}

	return val, err
}
