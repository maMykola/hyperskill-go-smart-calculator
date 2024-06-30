package calculator

import (
	"errors"
	"fmt"
	"strconv"
)

type Calculator struct {
	vars map[string]int
}

var (
	errNoExpression      = errors.New("No expression")
	errInvalidIdentifier = errors.New("Invalid identifier")
	errInvalidExpression = errors.New("Invalid expression")
	errInvalidAssignment = errors.New("Invalid assignment")
	errUnknownVariable   = errors.New("Unknown variable")
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
	tokens, err := buildPostfix(input)
	if err != nil {
		return 0, err
	}

	if len(tokens) == 0 {
		return 0, errNoExpression
	}

	var num int
	var stack = make([]int, 0, len(tokens)/2)

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, errInvalidExpression
			}
			l := len(stack)
			a, b := stack[l-2], stack[l-1]
			stack = stack[:l-2]
			stack = append(stack, doAction(a, b, token))
		default:
			if num, err = strconv.Atoi(token); err == nil {
				stack = append(stack, num)
			} else if num, err = c.get(token); err == nil {
				stack = append(stack, num)
			} else {
				return 0, err
			}
		}
	}

	if len(stack) != 1 {
		return 0, errInvalidExpression
	}

	return stack[0], nil
}

func doAction(a, b int, action string) int {
	switch action {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		panic("Unknown action")
	}
}
