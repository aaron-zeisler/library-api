package storage

import (
	"context"

	"github.com/aaron-zeisler/library-api/internal"
)

type BookStorage interface {
	GetBooks(ctx context.Context) ([]internal.Book, error)
	GetBookByID(ctx context.Context, bookID string) (internal.Book, error)
	CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error)
	UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error)
	DeleteBook(ctx context.Context, bookID string) error
}
