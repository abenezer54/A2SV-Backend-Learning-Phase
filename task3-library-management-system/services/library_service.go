package services

import (
	"errors"
	"fmt"

	"github.com/abenezer54/A2SV-Backend-Learning-Phase/tree/main/task3-library-management-system/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (lib *Library) AddBook(book models.Book) {
	lib.Books[book.ID] = book
}

func (lib *Library) RemoveBook(bookID int) {
	delete(lib.Books, bookID)
}

func (lib *Library) BorrowBook(bookID int, memberID int) error {
	_, ok := lib.Members[memberID]
	if !ok {
		return errors.New("invalid member id")
	}
	// have the copy of the value of the book
	book, ok := lib.Books[bookID]
	if !ok {
		return errors.New("invalid book id")
	}

	if lib.Books[bookID].Status != "Available" {
		return errors.New("book is not available now")
	}
	// update it's status
	book.Status = "Borrowed"
	lib.Books[bookID] = book // update to library
	return nil
}

func (lib *Library) ReturnBook(bookID int, memberID int) error {
	_, ok := lib.Members[memberID]
	if !ok {
		return errors.New("invalid member id")
	}
	book, ok := lib.Books[bookID]
	if !ok {
		return errors.New("invalid book id")
	}
	if lib.Books[bookID].Status != "Borrowed" {
		msg := fmt.Sprintf("book with id %v is not borrowed.", bookID)
		return errors.New(msg)
	}
	book.Status = "Available"
	lib.Books[bookID] = book
	return nil
}

func (lib *Library) ListAvailableBooks() []models.Book {
	listOfAvailableBooks := []models.Book{}
	for _, val := range lib.Books {
		if val.Status == "Available" {
			listOfAvailableBooks = append(listOfAvailableBooks, val)
		}
	}
	return listOfAvailableBooks
}

func (lib *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	_, ok := lib.Members[memberID]
	if !ok {
		msg := fmt.Sprintf("unable fo find member with id: %v.", memberID)
		return nil, errors.New(msg)
	}
	member := lib.Members[memberID]

	return member.BorrowedBooks, nil
}
