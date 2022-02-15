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

type UniqOptions struct {
	InputFileUsed  bool
	InputFileName  string
	OutputFileUsed bool
	OutputFileName string
	CduUsed        bool
	CduParam       string
	FNumber        int
	SNumber        int
	IUsed          bool
}

type parsingError struct {
	errorMsg string
}

func (err parsingError) Error() string {
	return err.errorMsg
}

func ParseArgs(args []string) (returnArgs UniqOptions, err error) {
	skipNextIter := false
	for index, param := range args {
		if skipNextIter {
			skipNextIter = false
			continue
		}
		switch param {
		case "-c":
			fallthrough
		case "-d":
			fallthrough
		case "-u":
			if returnArgs.CduUsed {
				err = parsingError{"-c -d -u params can't be used together"}
				return
			}
			returnArgs.CduParam = param
			returnArgs.CduUsed = true
		case "-f":
			if returnArgs.FNumber > 0 {
				err = parsingError{"-f already specified"}
				return
			}
			if len(args) == index+1 {
				err = parsingError{"-f specified, but num_fields not specified"}
				return
			}
			var convErr error
			returnArgs.FNumber, convErr = strconv.Atoi(args[index+1])
			skipNextIter = true
			if convErr != nil {
				err = parsingError{"-f specified, but num_fields not recognized as number"}
				return
			}
		case "-s":
			if returnArgs.SNumber > 0 {
				err = parsingError{"-s already specified"}
				return
			}
			if len(args) == index+1 {
				err = parsingError{"-s specified, but num_chars not specified"}
				return
			}
			var convErr error
			returnArgs.SNumber, convErr = strconv.Atoi(args[index+1])
			skipNextIter = true
			if convErr != nil {
				err = parsingError{"-s specified, but num_chars not recognized as number"}
				return
			}
		case "-i":
			if returnArgs.IUsed == true {
				err = parsingError{"-i already specified"}
				return
			}
			returnArgs.IUsed = true
		default:
			if returnArgs.InputFileUsed {
				err = parsingError{"input_file [output_file] already specified"}
				return
			}
			returnArgs.InputFileName = param
			returnArgs.InputFileUsed = true

			if len(args) == index+1 {
				continue
			}

			nextParam := args[index+1]
			if strings.HasPrefix(nextParam, "-") {
				continue
			}

			returnArgs.OutputFileName = nextParam
			returnArgs.OutputFileUsed = true
			skipNextIter = true
		}
	}
	return
}

func Uniq(inputLines []string, parsed UniqOptions) (returnValue []string) {
	var ouputStrings []outputData

	for _, processingString := range inputLines {
		compareString := processingString

		for i := 0; i < parsed.FNumber; i++ {
			index := strings.IndexByte(compareString, ' ')
			if index == -1 {
				compareString = ""
				break
			}
			compareString = compareString[index+1 : len(compareString)-1]
		}

		if parsed.SNumber > 0 {
			leftBound := parsed.SNumber
			if leftBound > len(compareString)-1 {
				leftBound = len(compareString) - 1
			}

			if leftBound == -1 {
				leftBound = 0
			}

			rightBound := len(compareString) - 1
			if rightBound == -1 {
				rightBound = 0
			}

			compareString = compareString[leftBound:rightBound]
		}

		lastScannedItem := &outputData{0, "", ""}
		firstIteration := false

		if len(ouputStrings) != 0 {
			lastScannedItem = &ouputStrings[len(ouputStrings)-1]
		} else {
			firstIteration = true
		}

		if !firstIteration && ((parsed.IUsed && strings.EqualFold(lastScannedItem.compareStr, compareString)) ||
			(!parsed.IUsed && lastScannedItem.compareStr == compareString)) {
			lastScannedItem.count++
		} else {
			ouputStrings = append(ouputStrings, outputData{1, processingString, compareString})
		}
	}

	if parsed.CduUsed {
		switch parsed.CduParam {
		case "-c":
			for _, item := range ouputStrings {
				returnValue = append(returnValue, strconv.Itoa(item.count)+" "+item.str)
			}
		case "-d":
			for _, item := range ouputStrings {
				if item.count > 1 {
					returnValue = append(returnValue, item.str)
				}
			}
		case "-u":
			for _, item := range ouputStrings {
				if item.count == 1 {
					returnValue = append(returnValue, item.str)
				}
			}
		}
		return
	}

	for _, item := range ouputStrings {
		returnValue = append(returnValue, item.str)
	}
	return
}
