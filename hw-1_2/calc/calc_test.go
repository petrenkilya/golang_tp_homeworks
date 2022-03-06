package calc

import (
	"fmt"
	"math"
	"testing"
)

func TestCalculatorTableDriven(t *testing.T) {
	var tests = []struct {
		in    string
		out   float64
		error bool
	}{
		//Example tests
		{"(1+2)-3", 0, false},
		{"(1+2)*3", 9, false},
		//Plus tests
		{"1+4+17", 22, false},
		{"-4+6", 2, false},
		{"-10+0+0", -10, false},
		//Minus tests
		{"1-15-4", -18, false},
		{"-4-2-5", -11, false},
		{"-9-0-0", -9, false},
		//Multiplication
		{"2*3*4", 24, false},
		{"-2*2*3", -12, false},
		{"2*15*0", 0, false},
		{"15*2*1", 30, false},
		//Division
		{"10/5", 2, false},
		{"-10/2", -5, false},
		{"10/0", math.Inf(1), false},
		{"-10/0", math.Inf(-1), false},
		//Brackets
		{"7+2*2", 11, false},
		{"(7+2)*2", 18, false},
		{"-(2+6)", -8, false},
		{"2*(2+3)-(2+6)", 2, false},
		{"(2-3*(2+7))*2", -50, false},
		//Spaces
		{"-2     +5*   4", 18, false},
		{"  4/2   +   16", 18, false},
		//Floats
		{"2.3+3.8", 6.1, false},
		{"4.2 *   2 -3", 5.4, false},
		//Errors
		{"(3*(2+4)", 0, true},
		{"-(У лукоморья дуб зеленый)", 0, true},
		{"", 0, true},
	}
	for _, tt := range tests {
		testName := tt.in
		t.Run(testName, func(t *testing.T) {
			result, err := Calculator(tt.in)
			if tt.error && err == nil {
				t.Error("Expression: '" + tt.in + "', expected error, got '" + fmt.Sprint(result) + "'")
				return
			}

			if result != tt.out {
				t.Error("Expression: '" + tt.in + "', wanted '" + fmt.Sprint(tt.out) + "', got '" + fmt.Sprint(result) + "'")
			}
		})
	}
}
