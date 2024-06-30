package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Calculator struct {
	vars map[string]int
}

var (
	errInvalidIdentifier = errors.New("Invalid identifier")
	errUnknownVariable   = errors.New("Unknown variable")
	errInvalidAssignment = errors.New("Invalid assignment")
)

func (c *Calculator) Run() {
	c.vars = make(map[string]int, 10)

	for {
		var value int
		var err error

		action := getAction()

		switch action.name {
		case exitAction:
			fmt.Println("Bye!")
			return
		case helpAction:
			c.help()
		case variableAction:
			if value, err = c.get(action.variable); err == nil {
				fmt.Println(value)
			}
		case assignAction:
			err = c.assign(action.variable, action.input)
		case calculateAction:
			if value, err = c.calc(action.input); err == nil {
				fmt.Println(value)
			}
		default:
			fmt.Println("Unknown commandType")
		}

		if err != nil && !errors.Is(err, errNoExpression) {
			fmt.Println(err)
		}
	}
}

func (c *Calculator) help() {
	fmt.Println("The program allow to calculate the sum or difference of numbers")
}

func (c *Calculator) get(name string) (int, error) {
	if !isVariable(name) {
		return 0, errInvalidIdentifier
	}

	if val, ok := c.vars[name]; ok {
		return val, nil
	}

	return 0, errUnknownVariable
}

func (c *Calculator) assign(name string, input string) error {
	if !isVariable(name) {
		return errInvalidIdentifier
	}

	val, err := c.calc(input)
	if err == nil {
		c.vars[name] = val
	} else if !errors.Is(err, errUnknownVariable) {
		err = errInvalidAssignment
	}

	return err
}

func (c *Calculator) calc(input string) (int, error) {
	tokens, err := getTokens(input)
	if err != nil {
		return 0, err
	}

	if len(tokens) == 0 {
		return 0, errNoExpression
	}

	var result int
	var op = PlusOperation

	for _, token := range tokens {
		switch token {
		case "+", "-":
			op.Update(operationType(token[0]))
		default:
			if op == NoOperation {
				return 0, errInvalidExpression
			}

			var num int
			var err error

			if r, _ := utf8.DecodeRuneInString(token); unicode.IsLetter(r) {
				num, err = c.get(token)
			} else if num, err = strconv.Atoi(token); err != nil {
				err = errInvalidExpression
			}

			if err != nil {
				return 0, err
			}

			result = doOperation(op, result, num)
			op = NoOperation
		}
	}

	if op != NoOperation {
		return 0, errInvalidExpression
	}

	return result, nil
}

func doOperation(op operationType, a, b int) int {
	switch op {
	case PlusOperation:
		return a + b
	case MinusOperation:
		return a - b
	default:
		panic("invalid operation")
	}
}
