package books

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"

	"github.com/aaron-zeisler/library-api/internal"
)

type service struct {
	db     booksDB
	logger *logrus.Logger
}

type booksDB interface {
	GetBooks(ctx context.Context) ([]internal.Book, error)
	GetBookByID(ctx context.Context, bookID string) (internal.Book, error)
	CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error)
	UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error)
	DeleteBook(ctx context.Context, bookID string) error
}

func NewService(db booksDB, logger *logrus.Logger) service {
	return service{
		db:     db,
		logger: logger,
	}
}

func (s service) GetBooks(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	books, err := s.db.GetBooks(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to retrieve books from the database")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to retrieve books from the database: %w", err).Error(),
		}, nil
	}

	responseBody, err := json.Marshal(books)
	if err != nil {
		s.logger.WithError(err).Error("failed to encode the books into an http response")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to encode the books into an http response: %w", err).Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func (s service) GetBookByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bookID := request.PathParameters["book_id"]

	book, err := s.db.GetBookByID(ctx, bookID)
	if err != nil {
		s.logger.WithError(err).WithField("book_id", bookID).Error("failed to retrieve the book from the database")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to retrieve a book from the database: %w", err).Error(),
		}, nil
	}

	responseBody, err := json.Marshal(book)
	if err != nil {
		s.logger.WithError(err).Error("failed to encode the book into an http response")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to encode the book into an http response: %w", err).Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func (s service) CreateBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var book internal.Book
	err := json.Unmarshal([]byte(request.Body), &book)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode the request body into a book object")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to decode the request body into a book object: %w", err).Error(),
		}, nil
	}

	newBook, err := s.db.CreateBook(ctx, book.Title, book.Author, book.ISBN, book.Description)
	if err != nil {
		s.logger.WithError(err).Error("failed to create a new book in the database")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to create a new book in the database: %w", err).Error(),
		}, nil
	}

	responseBody, err := json.Marshal(newBook)
	if err != nil {
		s.logger.WithError(err).Error("failed to encode the book into an http response")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to encode the book into an http response: %w", err).Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func (s service) UpdateBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bookID := request.PathParameters["book_id"]

	var book internal.Book
	err := json.Unmarshal([]byte(request.Body), &book)
	if err != nil {
		s.logger.WithError(err).Error("failed to decode the request body into a book object")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to decode the request body into a book object: %w", err).Error(),
		}, nil
	}

	updatedBook, err := s.db.UpdateBook(ctx, bookID, book.Title, book.Author, book.ISBN, book.Description)
	if err != nil {
		s.logger.WithError(err).WithField("book_id", bookID).Error("failed to update the book in the database")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to update the book in the database: %w", err).Error(),
		}, nil
	}

	responseBody, err := json.Marshal(updatedBook)
	if err != nil {
		s.logger.WithError(err).Error("failed to encode the book into an http response")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to encode the book into an http response: %w", err).Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func (s service) DeleteBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bookID := request.PathParameters["book_id"]

	err := s.db.DeleteBook(ctx, bookID)
	if err != nil {
		s.logger.WithError(err).WithField("book_id", bookID).Error("failed to delete the book from the database")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Errorf("failed to delete the book from the database: %w", err).Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
