package main

import (
	"flag"
	"fmt"
	"lineUtils/uniq"
	"lineUtils/utilities"
	"os"
)

func main() {
	opts, ok := createFlagsToUniqOptsFromArgs()
	if !ok {
		return
	}

	inputFile := os.Stdin
	outputFile := os.Stdout
	var err error

	if len(flag.Args()) > 0 {
		inputFile, err = os.Open(flag.Args()[0])
		if err != nil {
			fmt.Println("Can't open input file: " + flag.Args()[0])
		}
		return
	}

	if len(flag.Args()) > 1 {
		outputFile, err = os.Create(flag.Args()[1])
		if err != nil {
			fmt.Println("Can't create output file: " + flag.Args()[1])
		}
		return
	}

	inputLines, err := utilities.LinesRead(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputLines := uniq.Uniq(inputLines, opts)

	err = utilities.LinesWrite(outputLines, outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}
