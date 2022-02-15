package calc

import (
	"fmt"
	"math"
	"testing"
)

func TestCalculatorTableDriven(t *testing.T) {
	var tests = []struct {
		in  string
		out float64
	}{
		//Example tests
		{"(1+2)-3", 0},
		{"(1+2)*3", 9},
		//Plus tests
		{"1+4+17", 22},
		{"-4+6", 2},
		{"-10+0+0", -10},
		//Minus tests
		{"1-15-4", -18},
		{"-4-2-5", -11},
		{"-9-0-0", -9},
		//Multiplication
		{"2*3*4", 24},
		{"-2*2*3", -12},
		{"2*15*0", 0},
		{"15*2*1", 30},
		//Division
		{"10/5", 2},
		{"-10/2", -5},
		{"10/0", math.Inf(1)},
		{"-10/0", math.Inf(-1)},
		//Brackets
		{"7+2*2", 11},
		{"(7+2)*2", 18},
		{"-(2+6)", -8},
		{"2*(2+3)-(2+6)", 2},
		{"(2-3*(2+7))*2", -50},
		//Spaces
		{"-2     +5*   4", 18},
		{"  4/2   +   16", 18},
		//Floats
		{"2.3+3.8", 6.1},
		{"4.2 *   2 -3", 5.4},
	}
	for _, tt := range tests {
		testName := tt.in
		t.Run(testName, func(t *testing.T) {
			result, _ := Calculator(tt.in)
			if result != tt.out {
				t.Error("Expression: '" + tt.in + "', wanted '" + fmt.Sprint(tt.out) + "', got '" + fmt.Sprint(result) + "'")
			}
		})
	}
}
