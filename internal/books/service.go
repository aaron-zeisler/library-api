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

//func (s service) CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error) {
func (s service) CreateBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//return s.db.CreateBook(ctx, title, author, isbn, description)
	return events.APIGatewayProxyResponse{}, nil
}

//func (s service) UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error) {
func (s service) UpdateBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//return s.db.UpdateBook(ctx, bookID, title, author, isbn, description)
	return events.APIGatewayProxyResponse{}, nil
}

//func (s service) DeleteBook(ctx context.Context, bookID string) error {
func (s service) DeleteBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//return s.db.DeleteBook(ctx, bookID)
	return events.APIGatewayProxyResponse{}, nil
}
