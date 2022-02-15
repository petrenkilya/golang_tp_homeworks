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
		options     string
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
			"",
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
			"-c",
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
			"-d",
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
			"-u",
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
			"-i",
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
			"-f 1",
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
			"-s 1",
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
			"-f 1 -c",
		},
	}
	for _, testCase := range testsUniq {
		testName := testCase.options + " testing"
		t.Run(testName, func(t *testing.T) {
			parsedArgs, err := ParseArgs(strings.Split(testCase.options, " "))
			if err != nil {
				t.Error(err)
			}
			testingLines := Uniq(testCase.inputLines, parsedArgs)

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

func TestParseArgsTableDriven(t *testing.T) {
	var testsPasringArgs = []struct {
		in  string
		out UniqOptions
	}{{"-c", UniqOptions{CduUsed: true, CduParam: "-c"}},
		{"-d", UniqOptions{CduUsed: true, CduParam: "-d"}},
		{"-u", UniqOptions{CduUsed: true, CduParam: "-u"}},
		{"-f 32", UniqOptions{FNumber: 32}},
		{"-s 9", UniqOptions{SNumber: 9}},
		{"-i", UniqOptions{IUsed: true}},
		{"-c -i -f 4 -s 8", UniqOptions{CduUsed: true, CduParam: "-c", FNumber: 4, SNumber: 8, IUsed: true}},
		{"input.txt", UniqOptions{InputFileUsed: true, InputFileName: "input.txt"}},

		{"input.txt output.txt", UniqOptions{InputFileUsed: true, InputFileName: "input.txt",
			OutputFileUsed: true, OutputFileName: "output.txt"}},

		{"-u -i -f 4 -s 8 input.txt output.txt", UniqOptions{CduUsed: true, CduParam: "-u", FNumber: 4,
			SNumber: 8, IUsed: true, InputFileUsed: true, InputFileName: "input.txt",
			OutputFileUsed: true, OutputFileName: "output.txt"}},
	}

	for _, testCase := range testsPasringArgs {
		testName := "Testing: '" + testCase.in + "'\n"
		t.Run(testName, func(t *testing.T) {
			testingArgs, err := ParseArgs(strings.Split(testCase.in, " "))
			if err != nil {
				t.Error(err)
			}
			if testingArgs != testCase.out {
				t.Error("Result incorrect, got :\n" + fmt.Sprint(testingArgs) +
					"\n expected :\n" + fmt.Sprint(testCase.out))
			}
		})
	}
}
