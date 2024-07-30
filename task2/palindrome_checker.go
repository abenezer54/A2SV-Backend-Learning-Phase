package main

import "strings"

func isPalindrome(s string) bool{
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), "")
	s = strings.ToLower(s)
	n := len(s)
	for i := 0; i < n / 2; i++ {
		if s[i] != s[n - i - 1]{
			return false
		}
	}
	return true
}