package calc

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"strings"
)

type calculatorError struct {
	errStr string
}

func (cE calculatorError) Error() string {
	return cE.errStr
}

type operatorT struct {
	priority  int
	operator  uint8
	isOpening bool
	isClosing bool
}

func newOperator(operator uint8) (operatorT, bool) {
	var value operatorT
	value.operator = operator
	switch operator {
	case '+':
		fallthrough
	case '-':
		value.priority = 1
	case '*':
		fallthrough
	case '/':
		value.priority = 2
	case '(':
		value.isOpening = true
	case ')':
		value.isClosing = true
	default:
		return value, false
	}
	return value, true
}

func (op operatorT) calculate(a float64, b float64) float64 {
	switch op.operator {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		return a / b
	default:
		return 1
	}
}

func reversePolish(expr string) (string, error) {
	operatorStack := stack.New()
	outputStr := ""
	firstNumberPrinted := false

MainLoop:
	for len(expr) > 0 {
		indexOfData := strings.IndexAny(expr, "1234567890()+-*/")
		if indexOfData == -1 {
			break
		}
		expr = expr[indexOfData:]

		var bufferOperatorChar uint8
		_, err := fmt.Sscanf(expr, "%c", &bufferOperatorChar)
		if err != nil {
			return "", calculatorError{"Operator not supported"}
		}

		//processing operators
		bufferOperator, isOperator := newOperator(bufferOperatorChar)
		processAsOperator := true
		if isOperator && !firstNumberPrinted && bufferOperator.priority == 1 {
			var bufferNumber float64
			_, err = fmt.Sscanf(expr, "%f", &bufferNumber)
			if err == nil {
				processAsOperator = false
			}
		}

		if isOperator && processAsOperator {
			for bufferOperator.isClosing {
				if operatorStack.Len() == 0 {
					return "", calculatorError{"Number of brackets mismatched"}
				}
				stackOperator := operatorStack.Pop().(operatorT)
				if stackOperator.isOpening {
					expr = expr[1:]
					continue MainLoop
				}

				outputStr = outputStr + string(stackOperator.operator) + " "
			}

			for operatorStack.Len() > 0 {
				stackOperator := operatorStack.Peek().(operatorT)
				if stackOperator.priority < bufferOperator.priority || bufferOperator.isOpening {
					break
				}
				operatorStack.Pop()

				outputStr = outputStr + string(stackOperator.operator) + " "
			}

			if bufferOperator.isOpening {
				firstNumberPrinted = false
			}

			operatorStack.Push(bufferOperator)
			expr = expr[1:]
			continue
		}

		//processing number
		var bufferNumber float64
		_, err = fmt.Sscanf(expr, "%f", &bufferNumber)
		if err != nil {
			return "", calculatorError{"Bad math expression"}
		}
		bufferStr := fmt.Sprintf("%f", bufferNumber)
		if bufferNumber < 0 {
			expr = expr[1:]
		}
		firstNumberPrinted = true

		indexOfOperator := strings.IndexAny(expr, " +-*/()")
		if indexOfOperator == -1 {
			expr = expr[len(expr):]
		} else {
			expr = expr[indexOfOperator:]
		}

		outputStr = outputStr + bufferStr + " "
	}

	for operatorStack.Len() > 0 {
		stackOperator := operatorStack.Pop().(operatorT)
		if stackOperator.isOpening {
			return "", calculatorError{"Number of brackets mismatched"}
		}

		outputStr = outputStr + string(stackOperator.operator) + " "
	}

	if len(outputStr) > 0 && outputStr[len(outputStr)-1] == ' ' {
		outputStr = outputStr[:len(outputStr)-1]
	}

	return outputStr, nil
}

func Calculator(expr string) (result float64, er error) {
	reversePolishStr, err := reversePolish(expr)
	if err != nil {
		er = err
		return
	}
	splittedPolish := strings.Split(reversePolishStr, " ")

	numberStack := stack.New()
	for _, item := range splittedPolish {
		var num float64
		_, err = fmt.Sscanf(item, "%f", &num)
		if err == nil {
			numberStack.Push(num)
			continue
		}

		operator, isOperator := newOperator(item[0])
		if isOperator {
			if numberStack.Len() < 1 {
				er = calculatorError{"Bad number of operands"}
				return
			}
			b := numberStack.Pop().(float64)

			a := 0.0
			if numberStack.Len() > 0 {
				a = numberStack.Pop().(float64)
			}
			numberStack.Push(operator.calculate(a, b))
			continue
		}
		er = calculatorError{"Bad math expression"}
		return
	}

	if numberStack.Len() > 1 {
		er = calculatorError{"Bad number of operators"}
	}
	result = numberStack.Pop().(float64)
	return
}
