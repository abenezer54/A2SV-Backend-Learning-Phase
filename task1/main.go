package main

import (
	"fmt"
)

type Subject struct {
	name  string
	grade float64
}

func display(name string, data []Subject) {
	fmt.Printf("\nName: %s\n", name)
	fmt.Println("Subjects and Grades:")

	for _, subject := range data {
		fmt.Printf("Subject: %s, Grade: %.2f\n", subject.name, subject.grade)
	}

	average := calcAverage(data)
	fmt.Printf("\nAverage Grade: %.2f\n", average)
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

func main() {
	var (
		name       string
		numSubject int
	)

	fmt.Print("Enter your name: ")
	fmt.Scan(&name)

	fmt.Print("Enter the number of subjects: ")
	fmt.Scan(&numSubject)

	data := make([]Subject, numSubject)

	for idx := 1; idx <= numSubject; idx++ {
		sub := Subject{}

		fmt.Printf("Enter the name of %v subject: ", idx)
		fmt.Scan(&sub.name)

		fmt.Printf("Enter the grade of %v subject: ", idx)
		fmt.Scan(&sub.grade)
		data[idx-1] = sub
	}

	display(name, data)
}
