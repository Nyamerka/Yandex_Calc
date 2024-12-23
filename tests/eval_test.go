package tests

import (
	"Yandex_Calc/internal/eval"
	"math/big"
	"testing"
)

func TestEval_ValidExpressions(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"1+1", 2},
		{"2*2", 4},
		{"10/2", 5},
		{"3-1", 2},
	}

	for _, tt := range tests {
		result, err := eval.Eval(tt.expression)
		if err != nil {
			t.Errorf("unexpected error for expression %q: %v", tt.expression, err)
		}

		actual := eval.BigratToFloat(result)
		if actual != tt.expected {
			t.Errorf("for expression %q, expected %v, got %v", tt.expression, tt.expected, actual)
		}
	}
}

func TestEval_InvalidExpressions(t *testing.T) {
	tests := []string{
		"1//1",
		"abc",
		"1++",
		"",
	}

	for _, expr := range tests {
		_, err := eval.Eval(expr)
		if err == nil {
			t.Errorf("expected an error for expression %q, got nil", expr)
		} else if err != eval.ErrInvalidExpression {
			t.Errorf("expected ErrInvalidExpression for %q, got %v", expr, err)
		}
	}
}

func TestBigratToFloat(t *testing.T) {
	tests := []struct {
		bigrat   *big.Rat
		expected float64
	}{
		{big.NewRat(10, 2), 5},
		{big.NewRat(1, 3), 0.3333333333},
	}

	for _, tt := range tests {
		actual := eval.BigratToFloat(tt.bigrat)
		if actual != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, actual)
		}
	}
}
