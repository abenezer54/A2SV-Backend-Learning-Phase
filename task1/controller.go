package main

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

func getInteger(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if line == "" {
		return 0, errors.New("please enter a value")
	}

	line = strings.TrimSpace(line)
	number, err := strconv.Atoi(line)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func getFloor(reader *bufio.Reader) (float64, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0.0, err
	}
	if line == "" {
		return 0, errors.New("please enter a value")
	}
	line = strings.TrimSpace(line)
	number, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return 0.0, err
	}
	return number, err
}
