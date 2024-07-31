# Library Management System

This project implements a simple console-based library management system in Go, featuring a separation of concerns with controllers, services, and models.

## Overview

The library management system allows users to:

- Add books
- Remove books
- Borrow books
- Return books
- List available books
- List borrowed books by member

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/library-management-system.git
   ```

2. Navigate to the project directory:

   ```sh
   cd library-management-system
   ```

3. Initialize the project:

   ```sh
   go mod init your_project
   ```

4. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

- **Add a Book**: Allows you to add a new book to the library. The system will prompt for the book's ID, title, and author.

- **Remove a Book**: Allows you to remove a book from the library. The system will prompt for the book's ID.

- **Borrow a Book**: Allows a member to borrow a book. The system will prompt for the book's ID and the member's ID.

- **Return a Book**: Allows a member to return a borrowed book. The system will prompt for the book's ID and the member's ID.

- **List Available Books**: Displays a list of all books that are currently available in the library.

- **List Borrowed Books**: Displays a list of all books borrowed by a specific member. The system will prompt for the member's ID.

## Running the Application

To run the application, execute the following command:

```sh
go run main.go
```
