package main

import (
	"strings"
)

func countWords(s string) map[string]int{
	words := strings.Fields(s)
	count := make(map[string]int)

	for _, v := range words {
		count[v] += 1
	}
	return count
}
