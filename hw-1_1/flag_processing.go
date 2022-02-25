package main

import (
	"flag"
	"lineUtils/uniq"
	"log"
)

func createFlagsToUniqOptsFromArgs() (uniq.Options, bool) {
	opts := uniq.Options{}
	flag.BoolVar(&opts.IFlagUsed, "i", false, "Case insensitive compare")
	flag.BoolVar(&opts.CFlagUsed, "c", false, "Count repeated lines")
	flag.BoolVar(&opts.DFlagUsed, "d", false, "Print only repeated lines")
	flag.BoolVar(&opts.UFlagUsed, "u", false, "Print only not repeated lines")
	flag.IntVar(&opts.FFlagNumber, "f", 0, "Not use in compare first n fields")
	flag.IntVar(&opts.SFlagNumber, "s", 0, "Not use in compare first n chars")
	flag.Parse()

	if (opts.CFlagUsed && opts.DFlagUsed) || (opts.CFlagUsed && opts.UFlagUsed) || (opts.DFlagUsed && opts.UFlagUsed) {
		log.Fatalf("-c -d -u params can't be used together")
		flag.PrintDefaults()
		return uniq.Options{}, true
	}
	return opts, false
}
