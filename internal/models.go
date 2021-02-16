package internal

import "fmt"

type Book struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	ISBN        string     `json:"isbn"`
	Description string     `json:"description"`
	Status      BookStatus `json:"book_status"`
}

type BookStatus string

const (
	CheckedIn  BookStatus = "in"
	CheckedOut BookStatus = "out"
)

type ErrBookNotFound struct {
	BookID string
}

func (e ErrBookNotFound) Error() string {
	return fmt.Sprintf("The book with ID '%s' was not found", e.BookID)
}
