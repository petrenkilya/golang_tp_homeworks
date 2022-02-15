package main

import (
	"bufio"
	"calc/calc"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputStr := scanner.Text()

	result, err := calc.Calculator(inputStr)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println(result)
}
