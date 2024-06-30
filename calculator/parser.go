package calculator

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	errNoExpression      = errors.New("No expression")
	errInvalidExpression = errors.New("Invalid expression")
)

func getTokens(input string) ([]string, error) {
	var r rune

	tokens := make([]string, 0, 10)
	number := make([]rune, 0, 10)

	for weight, start := 0, 0; start < len(input); start += weight {
		r, weight = utf8.DecodeRuneInString(input[start:])
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			number = append(number, r)
			continue
		}

		if len(number) > 0 {
			tokens = append(tokens, string(number))
			number = number[:0]
		}

		if unicode.IsSpace(r) {
			continue
		}

		if r == '+' || r == '-' {
			tokens = append(tokens, string(r))
		} else {
			return nil, errInvalidExpression
		}
	}

	if len(number) > 0 {
		tokens = append(tokens, string(number))
	}

	return tokens, nil
}
