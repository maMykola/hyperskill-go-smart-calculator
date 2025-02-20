package calculator

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

var supportedOperators = []rune{'+', '-', '*', '/', '^'}

func buildPostfix(input string) ([]string, error) {
	var result []string
	var operators []rune
	var token strings.Builder
	var r rune
	var hasNumber bool

	input = simplifyOperations(input)
	for weight, start := 0, 0; start < len(input); start += weight {
		r, weight = utf8.DecodeRuneInString(input[start:])
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			token.WriteRune(r)
			continue
		}

		if token.Len() > 0 {
			hasNumber = true
			result = append(result, token.String())
			token.Reset()
		} else {
			hasNumber = false
		}

		if slices.Contains(supportedOperators, r) {
			ops := fetchLowerOperators(&operators, r)
			if len(ops) == 0 && (r == '-' || r == '+') && !hasNumber {
				result = append(result, "0")
			}
			result = append(result, ops...)
			operators = append(operators, r)
		} else if r == '(' {
			operators = append(operators, r)
		} else if r == ')' {
			if ops, err := fetchBeforeParenthesis(&operators); err == nil {
				result = append(result, ops...)
			} else {
				return nil, err
			}
		} else {
			return []string{}, errInvalidExpression
		}
	}

	// add last value if exists
	if token.Len() > 0 {
		result = append(result, token.String())
	}

	// check for unclosed parenthesis
	if slices.Contains(operators, '(') {
		return []string{}, errInvalidExpression
	}

	// append remaining operators
	slices.Reverse(operators)
	for _, r = range operators {
		result = append(result, string(r))
	}

	return result, nil
}

func fetchLowerOperators(operators *[]rune, r rune) []string {
	var result []string

	for i := len(*operators) - 1; i >= 0; i-- {
		op := (*operators)[i]
		if op == '(' || isPriorityHigher(r, op) {
			break
		} else {
			result = append(result, string(op))
			*operators = (*operators)[:i]
		}
	}

	return result
}

func isPriorityHigher(r rune, op rune) bool {
	switch r {
	case '^':
		return true
	case '*', '/':
		return op == '+' || op == '-'
	default:
		return false
	}
}

func fetchBeforeParenthesis(operators *[]rune) ([]string, error) {
	var result []string

	for i := len(*operators) - 1; i >= 0; i-- {
		op := (*operators)[i]
		*operators = (*operators)[:i]

		if op == '(' {
			return result, nil
		}

		result = append(result, string(op))
	}

	return nil, errInvalidExpression
}

// simplifyOperations will remove all spaces and simplify multiple + and - operations
func simplifyOperations(input string) string {
	simplified := input

	for {
		simplified = strings.ReplaceAll(input, " ", "")
		simplified = strings.ReplaceAll(simplified, "--", "+")
		simplified = strings.ReplaceAll(simplified, "++", "+")
		simplified = strings.ReplaceAll(simplified, "+-", "-")
		simplified = strings.ReplaceAll(simplified, "-+", "-")

		if simplified == input {
			return input
		}

		input = simplified
	}
}
