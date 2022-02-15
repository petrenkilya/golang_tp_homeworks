package main

import (
	"bufio"
	"fmt"
	"os"
	"uniqGo/uniq"
)

func main() {
	args := os.Args[1:]
	errorFile := os.Stderr

	parsed, err := uniq.ParseArgs(args)

	if err != nil {
		fmt.Fprintln(errorFile, err.Error())
		fmt.Fprintln(errorFile, "Usage "+os.Args[0]+" [-c | -d | -u] [-i] [-f num_fields] [-s num_chars] [input_file [output_file]]")
		return
	}

	inputFile := os.Stdin
	outputFile := os.Stdout

	if parsed.InputFileUsed {
		inputFile, err = os.Open(parsed.InputFileName)
		if err != nil {
			fmt.Fprintln(errorFile, "Error while opening input file: "+parsed.InputFileName)
			return
		}
		defer inputFile.Close()
	}

	if parsed.OutputFileUsed {
		outputFile, err = os.Create(parsed.OutputFileName)

		if err != nil {
			fmt.Fprintln(errorFile, "Error while opening output file: "+parsed.OutputFileName)
		}
		defer outputFile.Close()
	}

	scanner := bufio.NewScanner(inputFile)

	var scannedLines []string

	for scanner.Scan() {
		scannedLines = append(scannedLines, scanner.Text())
	}

	resultedLines := uniq.Uniq(scannedLines, parsed)

	for _, line := range resultedLines {
		fmt.Fprintln(outputFile, line)
	}
}
