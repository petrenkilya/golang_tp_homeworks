package uniq

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestUniqTableDriven(t *testing.T) {
	var testsUniq = []struct {
		inputLines  []string
		outputLines []string
		opts        Options
	}{
		{strings.Split(
			//InputLines
			"I love music.\n"+
				"I love music.\n"+
				"I love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.\n"+
				"Thanks.\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.",
			"\n"), strings.Split(
			//Output Lines
			"I love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"Thanks.\n"+
				"I love music of Kartik.",
			"\n"),
			//Args
			Options{},
		},

		{strings.Split(
			//InputLines
			"I love music.\n"+
				"I love music.\n"+
				"I love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.\n"+
				"Thanks.\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.",
			"\n"), strings.Split(
			//Output Lines
			"3 I love music.\n"+
				"1 \n"+
				"2 I love music of Kartik.\n"+
				"1 Thanks.\n"+
				"2 I love music of Kartik.",
			"\n"),
			//Args
			Options{CUsed: true},
		},

		{strings.Split(
			//InputLines
			"I love music.\n"+
				"I love music.\n"+
				"I love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.\n"+
				"Thanks.\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.",
			"\n"), strings.Split(
			//Output Lines
			"I love music.\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.",
			"\n"),
			//Args
			Options{DUsed: true},
		},

		{strings.Split(
			//InputLines
			"I love music.\n"+
				"I love music.\n"+
				"I love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.\n"+
				"Thanks.\n"+
				"I love music of Kartik.\n"+
				"I love music of Kartik.",
			"\n"), strings.Split(
			//Output Lines
			"\n"+
				"Thanks.",
			"\n"),
			//Args
			Options{UUsed: true},
		},

		{strings.Split(
			//InputLines
			"I LOVE MUSIC.\n"+
				"I love music.\n"+
				"I LoVe MuSiC.\n"+
				"\n"+
				"I love MuSIC of Kartik.\n"+
				"I love music of kartik.\n"+
				"Thanks.\n"+
				"I love music of kartik.\n"+
				"I love MuSIC of Kartik.",
			"\n"), strings.Split(
			//Output Lines
			"I LOVE MUSIC.\n"+
				"\n"+
				"I love MuSIC of Kartik.\n"+
				"Thanks.\n"+
				"I love music of kartik.",
			"\n"),
			//Args
			Options{IUsed: true},
		},

		{strings.Split(
			//InputLines
			"We love music.\n"+
				"I love music.\n"+
				"They love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"We love music of Kartik.\n"+
				"Thanks.",
			"\n"), strings.Split(
			//Output Lines
			"We love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"Thanks.",
			"\n"),
			//Args
			Options{FNumber: 1},
		},

		{strings.Split(
			//InputLines
			"I love music.\n"+
				"A love music.\n"+
				"C love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"We love music of Kartik.\n"+
				"Thanks.",
			"\n"), strings.Split(
			//Output Lines
			"I love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"We love music of Kartik.\n"+
				"Thanks.",
			"\n"),
			//Args
			Options{SNumber: 1},
		},

		{strings.Split(
			//InputLines
			"We love music.\n"+
				"I love music.\n"+
				"They love music.\n"+
				"\n"+
				"I love music of Kartik.\n"+
				"We love music of Kartik.\n"+
				"Thanks.",
			"\n"), strings.Split(
			//Output Lines
			"3 We love music.\n"+
				"1 \n"+
				"2 I love music of Kartik.\n"+
				"1 Thanks.",
			"\n"),
			//Args
			Options{FNumber: 1, CUsed: true},
		},
	}
	for _, testCase := range testsUniq {
		testName := fmt.Sprint(testCase.opts) + " testing"
		t.Run(testName, func(t *testing.T) {

			testingLines := Uniq(testCase.inputLines, testCase.opts)

			if len(testingLines) != len(testCase.outputLines) {
				t.Error("Number of output lines mismatched! Wanted: " + strconv.Itoa(len(testCase.outputLines)) + " , got: " + strconv.Itoa(len(testingLines)) + " \n")
				return
			}

			for index, expectedLine := range testCase.outputLines {
				if testingLines[index] != expectedLine {
					t.Error("Line " + strconv.Itoa(index) + " mismatched, got '" + testingLines[index] + "', expected '" + expectedLine + "'\n")
				}
			}
		})
	}
}
