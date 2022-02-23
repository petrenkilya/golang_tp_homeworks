package utilities

import (
	"bufio"
	"io"
)

func LinesWrite(input []string, ioWriter io.Writer) error {
	writer := bufio.NewWriter(ioWriter)
	for _, line := range input {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func LinesRead(ioReader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(ioReader)
	var scannedLines []string

	for scanner.Scan() {
		scannedLines = append(scannedLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return scannedLines, nil
}
