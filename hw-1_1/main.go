package main

import (
	"flag"
	"fmt"
	"lineUtils/uniq"
	"lineUtils/utilities"
	"os"
)

func main() {
	opts := uniq.Options{}

	flag.BoolVar(&opts.IUsed, "i", false, "Case insensitive compare")
	flag.BoolVar(&opts.CUsed, "c", false, "Count repeated lines")
	flag.BoolVar(&opts.DUsed, "d", false, "Print only repeated lines")
	flag.BoolVar(&opts.UUsed, "u", false, "Print only not repeated lines")
	flag.IntVar(&opts.FNumber, "f", 0, "Not use in compare first n fields")
	flag.IntVar(&opts.SNumber, "s", 0, "Not use in compare first n chars")
	flag.Parse()

	if (opts.CUsed && opts.DUsed) || (opts.CUsed && opts.UUsed) || (opts.DUsed && opts.UUsed) {
		fmt.Println("-c -d -u params can't be used together")
		flag.PrintDefaults()
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
