package main

import (
	"testing"
)

func TestCalcAverage(t *testing.T) {
	tests := []struct {
		name     string
		data     []Subject
		expected float64
	}{
		{"Single Subject", []Subject{{"DSA", 90}}, 90},
		{"Two Subjects", []Subject{{"DSA", 90}, {"OS", 80}}, 85},
		{"Multiple Subjects", []Subject{{"DSA", 90}, {"OS", 80}, {"SRE", 70}}, 80},
		{"No Subjects", []Subject{}, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calcAverage(test.data)
			if result != test.expected {
				t.Errorf("calcAverage(%v) = %f; want %f", test.data, result, test.expected)
			}
		})
	}
}
