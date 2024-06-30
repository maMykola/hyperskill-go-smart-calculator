package calculator

import (
	"bufio"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type commandType int

type actionType struct {
	name     commandType
	variable string
	input    string
}

const (
	unknownAction commandType = iota - 1
	exitAction
	helpAction
	variableAction
	assignAction
	calculateAction
)

var scanner = bufio.NewScanner(os.Stdin)

func getAction() actionType {
	scanner.Scan()
	data := strings.TrimSpace(scanner.Text())

	// predefined commands
	if strings.HasPrefix(data, "/") {
		return actionType{name: getCommand(data[1:])}
	}

	// just variable name
	if isVariable(data) {
		return actionType{
			name:     variableAction,
			variable: data,
		}
	}

	// assignment operation
	if i := strings.Index(data, "="); i >= 0 {
		return actionType{
			name:     assignAction,
			variable: strings.TrimSpace(data[:i]),
			input:    strings.TrimSpace(data[i+1:]),
		}
	}

	// calculation in all other cases
	return actionType{
		name:  calculateAction,
		input: data,
	}
}

func getCommand(data string) commandType {
	switch data {
	case "exit":
		return exitAction
	case "help":
		return helpAction
	default:
		return unknownAction
	}
}

func isVariable(data string) bool {
	if len(data) == 0 {
		return false
	}

	var r rune

	for weight, start := 0, 0; start < len(data); start += weight {
		r, weight = utf8.DecodeRuneInString(data[start:])
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
