package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Subject struct {
	name  string
	grade float64
}

var cont string = "yes"

func main() {
	for cont != "no" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("------------------------------------------")
		fmt.Println("Welcome to grade calculator")
		fmt.Println("------------------------------------------")
		fmt.Println("Enter your name: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ERROR!!!: ", err)
			continue
		}

		fmt.Println("Enter the number of the subjects")
		numSubject, err := getInteger(reader)
		if err != nil {
			fmt.Println("ERROR!!!: ", err)
			continue
		}

		data := make([]Subject, numSubject)

		for idx := 1; idx <= numSubject; idx++ {
			sub := Subject{}

			fmt.Printf("Enter the name of %v subject: \n", idx)
			sub.name, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("ERROR!!!: ", err)
				continue
			}

			fmt.Printf("Enter the grade of %v subject: \n", idx)
			sub.grade, err = getFloor(reader)
			if err != nil {
				fmt.Println("ERROR!!!: ", err)
				continue
			}
			msg, ok := validateGrade(sub.grade)
			for !ok {
				fmt.Println(msg, "||Enter a valid grade")
				sub.grade, err = getFloor(reader)
				if err != nil {
					fmt.Println("ERROR!!!: ", err)
					continue
				}
				_, ok = validateGrade(sub.grade)
			}
			data[idx-1] = sub
		}

		display(name, data)
		fmt.Println("------------------------------------------")
		fmt.Println("Do you want to continue? yes/no")
		fmt.Scan(&cont)
		cont = strings.ToLower(cont)
	}
}
