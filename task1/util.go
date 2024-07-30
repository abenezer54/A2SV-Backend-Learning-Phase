package main

import (
	"fmt"
)

func display(name string, data []Subject) {
	const (
		reset     = "\033[0m"
		cyan      = "\033[36m"
		yellow    = "\033[33m"
		green     = "\033[32m"
		magenta   = "\033[35m"
		bold      = "\033[1m"
		underline = "\033[4m"
	)

	fmt.Printf("\n%s%sName: %s%s\n", cyan, bold, name, reset)
	fmt.Printf("%s%sSubjects and Grades%s\n", yellow, underline, reset)
	fmt.Println("----------------------------------")

	for _, subject := range data {
		fmt.Printf("%sSubject: %-15s%s", green, subject.name, reset)
		fmt.Printf("%s| Grade: %.2f%s\n", magenta, subject.grade, reset)
	}

	fmt.Println("----------------------------------")
	average := calcAverage(data)
	fmt.Printf("\n%sAverage Grade: %.2f%s\n", magenta, average, reset)
}

func calcAverage(data []Subject) float64 {
	if len(data) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range data {
		sum += v.grade
	}

	average := sum / float64(len(data))
	return average
}
