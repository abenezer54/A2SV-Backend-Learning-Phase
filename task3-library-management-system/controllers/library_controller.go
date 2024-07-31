package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/abenezer54/A2SV-Backend-Learning-Phase/tree/main/task3-library-management-system/models"
	"github.com/abenezer54/A2SV-Backend-Learning-Phase/tree/main/task3-library-management-system/services"
)

type LibraryController struct {
	Library services.Library
}

func NewLibraryController(lib services.Library) *LibraryController {
	return &LibraryController{
		Library: lib,
	}
}

func (lc *LibraryController) AddBook() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter book ID: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

	fmt.Print("Enter book title: ")
	title, _ := reader.ReadString('\n')

	fmt.Print("Enter book author: ")
	author, _ := reader.ReadString('\n')

	book := models.NewBook(id, title, author)
	lc.Library.AddBook(*book)
	fmt.Println("Book added successfully!")
	fmt.Println("BOOK ID: ", id)
}

func (lc *LibraryController) RemoveBook() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter book ID to remove: ")
	idStr, _ := reader.ReadString('\n')
	id, _ := strconv.Atoi(idStr)
	lc.Library.RemoveBook(id)
}

func (lc *LibraryController) BorrowBook() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter book id: ")
	idStr, _ := reader.ReadString('\n')
	bookID, _ := strconv.Atoi(idStr)
	fmt.Println("Enter member id: ")
	idStr, _ = reader.ReadString('\n')
	memberID, _ := strconv.Atoi(idStr)

	err := lc.Library.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Unable to borrow book ||", err)
	}
}

func (lc *LibraryController) ReturnBook() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter book id: ")
	idStr, _ := reader.ReadString('\n')
	bookID, _ := strconv.Atoi(idStr)
	fmt.Println("Enter member id: ")
	idStr, _ = reader.ReadString('\n')
	memberID, _ := strconv.Atoi(idStr)

	err := lc.Library.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Unable to return book ||", err)
	}
}

func (lc *LibraryController) ListAvailableBooks() {
	books := lc.Library.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}

	fmt.Println("Available Books:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%-5s %-25s %-20s\n", "ID", "Title", "Author")
	fmt.Println(strings.Repeat("-", 50))
	for _, book := range books {
		fmt.Printf("%-5d %-25s %-20s\n", book.ID, book.Title, book.Author)
	}
	fmt.Println(strings.Repeat("-", 50))
}

func (lc *LibraryController) ListBorrowedBooks() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter member id: ")
	idStr, _ := reader.ReadString('\n')
	memberID, _ := strconv.Atoi(idStr)

	books, err := lc.Library.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(books) == 0 {
		fmt.Println("No available books.")
		return
	}

	fmt.Println("Available Books:")
	fmt.Println("ID\tTitle\tAuthor")
	for _, book := range books {
		fmt.Printf("%d\t%s\t%s\n", book.ID, book.Title, book.Author)
	}
}
