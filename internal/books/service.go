package books

import (
	"context"

	"github.com/aaron-zeisler/library-api/internal"
)

type service struct {
	db booksDB
}

type booksDB interface {
	GetBooks(ctx context.Context) ([]internal.Book, error)
	GetBookByID(ctx context.Context, bookID string) (internal.Book, error)
	CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error)
	UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error)
	DeleteBook(ctx context.Context, bookID string) error
}

func NewService(ctx context.Context, db booksDB) service {
	return service{
		db: db,
	}
}

func (s service) GetBooks(ctx context.Context) ([]internal.Book, error) {
	return s.db.GetBooks(ctx)
}

func (s service) GetBookByID(ctx context.Context, bookID string) (internal.Book, error) {
	return s.db.GetBookByID(ctx, bookID)
}

func (s service) CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error) {
	return s.db.CreateBook(ctx, title, author, isbn, description)
}

func (s service) UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error) {
	return s.db.UpdateBook(ctx, bookID, title, author, isbn, description)
}

func (s service) DeleteBook(ctx context.Context, bookID string) error {
	return s.db.DeleteBook(ctx, bookID)
}
