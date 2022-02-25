package main

import (
	"bufio"
	"calc/calc"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputStr := scanner.Text()

	result, err := calc.Calculator(inputStr)

	if err != nil {
		log.Println(err)
		log.Printf("%w", err)
		return
	}

	fmt.Println(result)
}
