package calc

import (
	"fmt"
	"strings"
)

type Stack struct {
	data []interface{}
	head int
}

func NewStack() *Stack {
	return &Stack{head: 0, data: []interface{}{}}
}

func (st *Stack) Push(elem interface{}) {
	if st.head == len(st.data) {
		st.data = append(st.data, len(st.data)*2)
	}
	st.data[st.head] = elem
	st.head++
}

func (st *Stack) Pop() interface{} {
	if st.head == 0 {
		panic("empty stack pop")
	}
	st.head--
	return st.data[st.head]
}

func (st *Stack) Peek() interface{} {
	if st.head == 0 {
		panic("empty stack peek")
	}
	return st.data[st.head-1]
}

func (st *Stack) Len() int {
	return st.head
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

func cutStrToAny(strToCut *string, charsToFind string) int {
	indexOfFoundData := strings.IndexAny(*strToCut, charsToFind)
	if indexOfFoundData == -1 {
		*strToCut = (*strToCut)[len(*strToCut):]
		return indexOfFoundData
	}
	*strToCut = (*strToCut)[indexOfFoundData:]
	return 0
}

func finishPolishNotation(operatorStack *Stack, inputStr string) (string, error) {
	for operatorStack.Len() > 0 {
		stackOperator := operatorStack.Pop().(operatorT)
		if stackOperator.isOpening {
			return "", fmt.Errorf("number of brackets mismatched")
		}

		inputStr = inputStr + string(stackOperator.operator) + " "
	}

	if len(inputStr) > 0 && (inputStr)[len(inputStr)-1] == ' ' {
		inputStr = (inputStr)[:len(inputStr)-1]
	}
	return inputStr, nil
}

func reversePolishProcessLowerPriorityOperator(operatorStack *Stack, bufferOperator operatorT, outputStr *string) {
	for operatorStack.Len() > 0 {
		stackOperator := operatorStack.Peek().(operatorT)
		if stackOperator.priority < bufferOperator.priority || bufferOperator.isOpening {
			break
		}
		operatorStack.Pop()

		*outputStr = *outputStr + string(stackOperator.operator) + " "
	}
}

func reversePolishProcessClosing(operatorStack *Stack, bufferOperator operatorT, outputStr *string) (bool, error) {
	for bufferOperator.isClosing {
		if operatorStack.Len() == 0 {
			return false, fmt.Errorf("number of brackets mismatched")
		}
		stackOperator := operatorStack.Pop().(operatorT)
		if stackOperator.isOpening {
			return true, nil
		}

		*outputStr = *outputStr + string(stackOperator.operator) + " "
	}
	return false, nil
}

func reversePolishNotation(expr string) (string, error) {
	operatorStack := NewStack()
	outputStr := ""
	firstNumberPrinted := false

	for len(expr) > 0 {
		if cutStrToAny(&expr, "1234567890()+-*/") == -1 {
			break
		}

		var bufferOperatorChar uint8
		_, err := fmt.Sscanf(expr, "%c", &bufferOperatorChar)
		if err != nil {
			return "", fmt.Errorf("operator not supported")
		}

		bufferOperator, isOperator := newOperator(bufferOperatorChar)
		processAsOperator := true

		//try to process as number
		if isOperator && !firstNumberPrinted && bufferOperator.priority == 1 {
			var bufferNumber float64
			_, err = fmt.Sscanf(expr, "%f", &bufferNumber)
			if err == nil {
				processAsOperator = false
			}
		}

		//processing operators
		if isOperator && processAsOperator {
			processed, err := reversePolishProcessClosing(operatorStack, bufferOperator, &outputStr)
			if err != nil {
				return "", err
			}
			if processed {
				expr = expr[1:]
				continue
			}

			reversePolishProcessLowerPriorityOperator(operatorStack, bufferOperator, &outputStr)

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
			return "", fmt.Errorf("bad math expression")
		}
		bufferStr := fmt.Sprintf("%f", bufferNumber)
		if bufferNumber < 0 {
			expr = expr[1:]
		}
		firstNumberPrinted = true

		cutStrToAny(&expr, " +-*/()")

		outputStr = outputStr + bufferStr + " "
	}

	return finishPolishNotation(operatorStack, outputStr)
}

func Calculator(expr string) (float64, error) {
	reversePolishStr, err := reversePolishNotation(expr)
	if err != nil {
		return 0, err
	}

	if len(reversePolishStr) == 0 {
		return 0, fmt.Errorf("bad math expression")
	}
	splittedPolish := strings.Split(reversePolishStr, " ")

	numberStack := NewStack()
	for _, item := range splittedPolish {
		var num float64
		_, err = fmt.Sscanf(item, "%f", &num)
		if err == nil {
			numberStack.Push(num)
			continue
		}

		operator, isOperator := newOperator(item[0])
		if !isOperator {
			return 0, fmt.Errorf("bad math expression")
		}

		if numberStack.Len() < 1 {
			return 0, fmt.Errorf("bad number of operands")
		}
		b := numberStack.Pop().(float64)

		a := 0.0
		if numberStack.Len() > 0 {
			a = numberStack.Pop().(float64)
		}
		numberStack.Push(operator.calculate(a, b))
		continue
	}

	if numberStack.Len() > 1 {
		return 0, fmt.Errorf("bad number of operators")
	}

	return numberStack.Pop().(float64), nil
}
