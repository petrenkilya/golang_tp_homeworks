package uniq

import (
	"strconv"
	"strings"
)

type outputData struct {
	count      int
	str        string
	compareStr string
}

type Options struct {
	CUsed   bool
	DUsed   bool
	UUsed   bool
	FNumber int
	SNumber int
	IUsed   bool
}

func formResult(outputStrings []outputData, opts Options) (returnValue []string) {
	if opts.CUsed {
		returnValue = make([]string, len(outputStrings), len(outputStrings))
		for index, item := range outputStrings {
			returnValue[index] = strconv.Itoa(item.count) + " " + item.str
		}
	} else if opts.DUsed {
		for _, item := range outputStrings {
			if item.count > 1 {
				returnValue = append(returnValue, item.str)
			}
		}
	} else if opts.UUsed {
		for _, item := range outputStrings {
			if item.count == 1 {
				returnValue = append(returnValue, item.str)
			}
		}
	} else {
		returnValue = make([]string, len(outputStrings), len(outputStrings))
		for index, item := range outputStrings {
			returnValue[index] = item.str
		}
	}
	return
}

func getCompareString(inputStr string, opts Options) string {
	for i := 0; i < opts.FNumber; i++ {
		index := strings.IndexByte(inputStr, ' ')
		if index == -1 {
			inputStr = ""
			break
		}
		inputStr = inputStr[index+1 : len(inputStr)-1]
	}

	if opts.SNumber > 0 {
		leftBound := opts.SNumber
		if leftBound > len(inputStr)-1 {
			leftBound = len(inputStr) - 1
		}

		if leftBound == -1 {
			leftBound = 0
		}

		rightBound := len(inputStr) - 1
		if rightBound == -1 {
			rightBound = 0
		}

		inputStr = inputStr[leftBound:rightBound]
	}

	return inputStr
}

func Uniq(inputLines []string, opts Options) (returnValue []string) {
	var ouputStrings []outputData

	for _, processingString := range inputLines {
		compareString := getCompareString(processingString, opts)

		lastScannedItem := &outputData{0, "", ""}
		firstIteration := false

		if len(ouputStrings) != 0 {
			lastScannedItem = &ouputStrings[len(ouputStrings)-1]
		} else {
			firstIteration = true
		}

		if !firstIteration && ((opts.IUsed && strings.EqualFold(lastScannedItem.compareStr, compareString)) ||
			(!opts.IUsed && lastScannedItem.compareStr == compareString)) {
			lastScannedItem.count++
		} else {
			ouputStrings = append(ouputStrings, outputData{1, processingString, compareString})
		}
	}

	returnValue = formResult(ouputStrings, opts)

	return
}
