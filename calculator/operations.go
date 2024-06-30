package calculator

type operationType byte

const (
	NoOperation    operationType = 0
	PlusOperation  operationType = '+'
	MinusOperation operationType = '-'
)

func (op *operationType) Update(action operationType) {
	if *op == NoOperation || *op == PlusOperation {
		*op = action
	} else if action == MinusOperation {
		*op = PlusOperation
	} else {
		*op = action
	}
}

func isOperator(op string) bool {
	switch op {
	case "+", "-", "*", "/", "(", ")":
		return true
	default:
		return false
	}
}
