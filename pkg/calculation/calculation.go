package calculation

import (
	"strconv"
	"strings"
	"unicode"
)

func stringToFloat64(str string) (float64, error) {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, ErrInvalidExpression
	}
	return res, nil
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}

func precedence(op rune) int {
	if op == '+' || op == '-' {
		return 1
	}
	if op == '*' || op == '/' {
		return 2
	}
	return 0
}

func applyOp(a, b float64, op rune) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	}
	return 0, ErrInvalidExpression
}
func Calc(expression string) (float64, error) {
	var values []float64
	var ops []rune

	expression = strings.ReplaceAll(expression, " ", "")

	for i := 0; i < len(expression); {
		char := rune(expression[i])

		if unicode.IsDigit(char) || char == '.' {
			var buffer strings.Builder
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				buffer.WriteByte(expression[i])
				i++
			}
			value, err := stringToFloat64(buffer.String())
			if err != nil {
				return 0, err
			}
			values = append(values, value)
		} else if char == '(' {
			ops = append(ops, char)
			i++
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(values) < 2 {
					return 0, ErrInvalidExpression
				}
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				val2 := values[len(values)-1]
				values = values[:len(values)-1]

				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				result, err := applyOp(val1, val2, op)
				if err != nil {
					return 0, err
				}
				values = append(values, result)
			}
			if len(ops) == 0 || ops[len(ops)-1] != '(' {
				return 0, ErrInvalidExpression
			}
			ops = ops[:len(ops)-1]
			i++
		} else if isSign(char) {
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(char) {
				if len(values) < 2 {
					return 0, ErrInvalidExpression
				}
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				val2 := values[len(values)-1]
				values = values[:len(values)-1]

				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				result, err := applyOp(val1, val2, op)
				if err != nil {
					return 0, err
				}
				values = append(values, result)
			}
			ops = append(ops, char)
			i++
		} else {
			return 0, ErrInvalidExpression
		}
	}

	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, ErrInvalidExpression
		}
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		val2 := values[len(values)-1]
		values = values[:len(values)-1]

		val1 := values[len(values)-1]
		values = values[:len(values)-1]

		result, err := applyOp(val1, val2, op)
		if err != nil {
			return 0, err
		}
		values = append(values, result)
	}

	if len(values) != 1 {
		return 0, ErrInvalidExpression
	}

	return values[0], nil
}
