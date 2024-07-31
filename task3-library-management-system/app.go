package main

import (
	"fmt"

	"github.com/abenezer54/A2SV-Backend-Learning-Phase/tree/main/task3-library-management-system/controllers"
	"github.com/abenezer54/A2SV-Backend-Learning-Phase/tree/main/task3-library-management-system/services"
)

func RunApp() {
	library := services.NewLibrary()
	controller := controllers.NewLibraryController(*library)

	for {
		fmt.Println("Library Management System")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			controller.AddBook()
		case 2:
			controller.RemoveBook()
		case 3:
			controller.BorrowBook()
		case 4:
			controller.ReturnBook()
		case 5:
			controller.ListAvailableBooks()
		case 6:
			controller.ListBorrowedBooks()
		case 7:
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
