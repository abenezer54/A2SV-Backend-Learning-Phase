package main

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

func getLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if err != nil {
		return "", err
	}
	if len(line) == 0 {
		return "", errors.New("please enter a value")
	}
	return "", nil
}

func getInteger(reader *bufio.Reader) (int, error) {
	line, err := getLine(reader)
	if err != nil {
		return 0, err
	}

	line = strings.TrimSpace(line)
	number, err := strconv.Atoi(line)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func getFloat(reader *bufio.Reader) (float64, error) {
	line, err := getLine(reader)
	if err != nil {
		return 0.0, err
	}

	line = strings.TrimSpace(line)
	number, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return 0.0, err
	}
	return number, err
}
