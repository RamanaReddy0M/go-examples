package main

import (
	"bufio"
	"os"
)

func ReadLines(filePath string) ([]string, error) {

	lines := make([]string, 0)

	fileReader, error := os.Open(filePath)
	//check for erros
	if error != nil {
		return nil, error
	}

	defer fileReader.Close()

	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
