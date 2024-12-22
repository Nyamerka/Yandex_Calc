package eval

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var fp_rx = regexp.MustCompile(`(\d+(?:\.\d+)?)`)
var functions_rx = regexp.MustCompile(`(sin|cos|tan|ln|arcsin|arccos|arctan|sqrt)`)
var symbols_rx *regexp.Regexp // TODO used as a mutable global variable!!
var unary_minus_rx = regexp.MustCompile(`((?:^|[-+^%*/<>!=(])\s*)-`)
var whitespace_rx = regexp.MustCompile(`\s+`)

var ErrInvalidExpression = errors.New("invalid expression")

var symbolTable map[string]string // TODO used as a mutable global variable!!

var operators = []string{"-", "+", "*", "/", "<", ">", "@", "^", "**", "%", "!=", "==", ">=", "<="}

func prec(op string) (result int) {
	if op == "-" || op == "+" {
		result = 1
	} else if op == "*" || op == "/" {
		result = 2
	} else if op == "^" || op == "%" || op == "**" {
		result = 3
	} else if op == "@" {
		result = 4
	} else if functions_rx.MatchString(op) {
		result = 5
	} else {
		result = 0
	}
	return
}

func opGTE(op1, op2 string) bool {
	return prec(op1) >= prec(op2)
}

func isFunction(token string) bool {
	return functions_rx.MatchString(token)
}

func isOperator(token string) bool {
	for _, v := range operators {
		if v == token {
			return true
		}
	}
	return false
}

func isOperand(token string) bool {
	return fp_rx.MatchString(token)
}

func isSymbol(token string) bool {
	for k := range symbolTable {
		if k == token {
			return true
		}
	}
	return false
}

func convert2postfix(tokens []string) []string {
	var stack Stack
	var result []string
	for _, token := range tokens {

		stackString := fmt.Sprint(stack)
		stackString += ""

		if isOperator(token) {
		OPERATOR:
			for {
				top, err := stack.Top()
				if err == nil && top != "(" {
					if opGTE(top.(string), token) {
						pop, _ := stack.Pop()
						result = append(result, pop.(string))
						continue
					} else {
						break OPERATOR
					}
				}
				break OPERATOR
			}
			stack.Push(token)

		} else if isFunction(token) {
		FUNCTION:
			for {
				top, err := stack.Top()
				if err == nil && top != "(" {
					if opGTE(top.(string), token) {
						pop, _ := stack.Pop()
						result = append(result, pop.(string))
					}
				} else {
					break FUNCTION
				}
				break FUNCTION
			}
			stack.Push(token)

		} else if token == "(" {
			stack.Push(token)

		} else if token == ")" {
		PAREN:
			for {
				top, err := stack.Top()
				if err == nil && top != "(" {
					pop, _ := stack.Pop()
					result = append(result, pop.(string))
				} else {
					stack.Pop()
					break PAREN
				}
			}

		} else if isOperand(token) {
			result = append(result, token)

		} else if isSymbol(token) {
			result = append(result, symbolTable[token])

		} else {
			result = append(result, token)
		}

	}

	for !stack.IsEmpty() {
		pop, _ := stack.Pop()
		result = append(result, pop.(string))
	}

	return result
}

func evaluatePostfix(postfix []string) (*big.Rat, error) {
	var stack Stack
	result := new(big.Rat)
	for _, token := range postfix {

		stackString := fmt.Sprint(stack)
		stackString += ""

		if isOperand(token) {
			bigrat := new(big.Rat)
			if _, err := fmt.Sscan(token, bigrat); err != nil {
				return nil, fmt.Errorf("unable to scan %s", token)
			}
			stack.Push(bigrat)

		} else if isOperator(token) {

			op2, err2 := stack.Pop()
			if err2 != nil {
				return nil, err2
			}

			var op1 interface{}
			if token != "@" {
				var err1 error
				if op1, err1 = stack.Pop(); err1 != nil {
					return nil, err1
				}
			}

			dummy := new(big.Rat)
			switch token {
			case "**", "^":
				float1 := BigratToFloat(op1.(*big.Rat))
				float2 := BigratToFloat(op2.(*big.Rat))
				float_result := math.Pow(float1, float2)
				stack.Push(FloatToBigrat(float_result))
			case "%":
				float1 := BigratToFloat(op1.(*big.Rat))
				float2 := BigratToFloat(op2.(*big.Rat))
				float_result := math.Mod(float1, float2)
				stack.Push(FloatToBigrat(float_result))
			case "*":
				result := dummy.Mul(op1.(*big.Rat), op2.(*big.Rat))
				stack.Push(result)
			case "/":
				result := dummy.Quo(op1.(*big.Rat), op2.(*big.Rat))
				stack.Push(result)
			case "+":
				result = dummy.Add(op1.(*big.Rat), op2.(*big.Rat))
				stack.Push(result)
			case "-":
				result = dummy.Sub(op1.(*big.Rat), op2.(*big.Rat))
				stack.Push(result)
			case "<":
				if op1.(*big.Rat).Cmp(op2.(*big.Rat)) <= -1 {
					stack.Push(big.NewRat(1, 1))
				} else {
					stack.Push(new(big.Rat))
				}
			case "<=":
				if op1.(*big.Rat).Cmp(op2.(*big.Rat)) <= 0 {
					stack.Push(big.NewRat(1, 1))
				} else {
					stack.Push(new(big.Rat))
				}
			case ">":
				if op1.(*big.Rat).Cmp(op2.(*big.Rat)) >= 1 {
					stack.Push(big.NewRat(1, 1))
				} else {
					stack.Push(new(big.Rat))
				}
			case ">=":
				if op1.(*big.Rat).Cmp(op2.(*big.Rat)) >= 0 {
					stack.Push(big.NewRat(1, 1))
				} else {
					stack.Push(new(big.Rat))
				}
			case "==":
				if op1.(*big.Rat).Cmp(op2.(*big.Rat)) == 0 {
					stack.Push(big.NewRat(1, 1))
				} else {
					stack.Push(new(big.Rat))
				}
			case "!=":
				if op1.(*big.Rat).Cmp(op2.(*big.Rat)) == 0 {
					stack.Push(new(big.Rat))
				} else {
					stack.Push(big.NewRat(1, 1))
				}
			case "@":
				result := dummy.Mul(big.NewRat(-1, 1), op2.(*big.Rat))
				stack.Push(result)
			}
		} else if isFunction(token) {
			op2, err := stack.Pop()
			if err != nil {
				return nil, err
			}
			switch token {
			case "sin":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Sin(float_result)))
			case "cos":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Cos(float_result)))
			case "tan":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Tan(float_result)))
			case "arcsin":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Asin(float_result)))
			case "arccos":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Acos(float_result)))
			case "arctan":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Atan(float_result)))
			case "ln":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Log(float_result)))
			case "sqrt":
				float_result := BigratToFloat(op2.(*big.Rat))
				stack.Push(FloatToBigrat(math.Sqrt(float_result)))
			}
		} else {
			return nil, ErrInvalidExpression
		}
	}

	retval, err := stack.Pop()
	if err != nil {
		return nil, err
	}
	return retval.(*big.Rat), nil
}

func Tokenise(expr string) []string {

	spaced := unary_minus_rx.ReplaceAllString(expr, "$1 @")
	spaced = fp_rx.ReplaceAllString(spaced, " ${1} ")
	spaced = functions_rx.ReplaceAllString(spaced, " ${1} ")

	if symbols_rx != nil {
		spaced = symbols_rx.ReplaceAllString(spaced, " ${1} ")
	}

	symbols := []string{"(", ")"}
	for _, symbol := range symbols {
		spaced = strings.Replace(spaced, symbol, fmt.Sprintf(" %s ", symbol), -1)
	}

	stripped := whitespace_rx.ReplaceAllString(strings.TrimSpace(spaced), "|")
	result := strings.Split(stripped, "|")
	return result
}

func Eval(expr string) (result *big.Rat, err error) {
	defer func() {
		if e := recover(); e != nil {
			result = nil
			err = ErrInvalidExpression
		}
	}()

	tokens := Tokenise(expr)
	postfix := convert2postfix(tokens)
	return evaluatePostfix(postfix)
}

func EvalWithVariables(expr string, variables map[string]string) (result *big.Rat, err error) {
	symbolTable = variables
	s := ""
	for k := range symbolTable {
		s += k
	}
	symbols_rx = regexp.MustCompile(fmt.Sprintf("(%s)", s))
	return Eval(expr)
}

func BigratToInt(bigrat *big.Rat) (int64, error) {
	float_string := bigrat.FloatString(0)
	return strconv.ParseInt(float_string, 10, 64)
}

func BigratToBigint(bigrat *big.Rat) *big.Int {
	int_string := bigrat.FloatString(0)
	bigint := new(big.Int)
	fmt.Sscan(int_string, bigint)
	return bigint
}

func BigratToFloat(bigrat *big.Rat) float64 {
	float_string := bigrat.FloatString(10)
	float, _ := strconv.ParseFloat(float_string, 64)
	return float
}

func FloatToBigrat(float float64) *big.Rat {
	float_string := fmt.Sprintf("%g", float)
	bigrat := new(big.Rat)
	fmt.Sscan(float_string, bigrat)
	return bigrat
}
