package books

import (
	"context"
	"encoding/json"
	"errors"
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
	UpdateBook(ctx context.Context, bookID string, book internal.Book) (internal.Book, error)
	DeleteBook(ctx context.Context, bookID string) error
}

func NewService(db booksDB, opts ...ServiceOption) service {
	s := service{
		db:     db,
		logger: logrus.New(),
	}

	for _, opt := range opts {
		s = opt(s)
	}

	return s
}

type ServiceOption func(s service) service

func WithLogger(logger *logrus.Logger) ServiceOption {
	return func(s service) service {
		s.logger = logger
		return s
	}
}

func (s service) GetBooks(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	books, err := s.db.GetBooks(ctx)
	if err != nil {
		return s.logAndReturnError(err, "failed to retrieve books from the database", http.StatusInternalServerError, logrus.Fields{})
	}

	responseBody, err := json.Marshal(books)
	if err != nil {
		return s.logAndReturnError(err, "failed to encode the books into an http response", http.StatusInternalServerError, logrus.Fields{})
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
		statusCode := http.StatusInternalServerError
		if errors.As(err, &internal.ErrBookNotFound{}) {
			statusCode = http.StatusNotFound
		}

		return s.logAndReturnError(err, "failed to retrieve the book from the database", statusCode, logrus.Fields{"book_id": bookID})
	}

	responseBody, err := json.Marshal(book)
	if err != nil {
		return s.logAndReturnError(err, "failed to encode the book into an http response", http.StatusInternalServerError, logrus.Fields{})
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
		return s.logAndReturnError(err, "failed to decode the request body into a book object", http.StatusBadRequest, logrus.Fields{})
	}

	newBook, err := s.db.CreateBook(ctx, book.Title, book.Author, book.ISBN, book.Description)
	if err != nil {
		return s.logAndReturnError(err, "failed to create a new book in the database", http.StatusInternalServerError, logrus.Fields{})
	}

	responseBody, err := json.Marshal(newBook)
	if err != nil {
		return s.logAndReturnError(err, "failed to encode the book into an http response", http.StatusInternalServerError, logrus.Fields{})
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
		return s.logAndReturnError(err, "failed to decode the request body into a book object", http.StatusBadRequest, logrus.Fields{})
	}

	updatedBook, err := s.db.UpdateBook(ctx, bookID, book)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &internal.ErrBookNotFound{}) {
			statusCode = http.StatusNotFound
		}

		return s.logAndReturnError(err, "failed to update the book in the database", statusCode, logrus.Fields{"book_id": bookID})
	}

	responseBody, err := json.Marshal(updatedBook)
	if err != nil {
		return s.logAndReturnError(err, "failed to encode the book into an http response", http.StatusInternalServerError, logrus.Fields{})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func (s service) DeleteBook(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bookID := request.PathParameters["book_id"]

	err := s.db.DeleteBook(ctx, bookID)
	if err != nil && !errors.As(err, &internal.ErrBookNotFound{}) { // 'Book not found' doesn't cause a 404 for the DELETE action
		return s.logAndReturnError(err, "failed to delete the book from the database", http.StatusInternalServerError, logrus.Fields{"book_id": bookID})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func (s service) CheckOut(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return s.updateStatus(ctx, request, internal.CheckedOut)
}

func (s service) CheckIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return s.updateStatus(ctx, request, internal.CheckedIn)
}

func (s service) updateStatus(ctx context.Context, request events.APIGatewayProxyRequest, newStatus internal.BookStatus) (events.APIGatewayProxyResponse, error) {
	bookID := request.PathParameters["book_id"]

	// Retrieve the book
	book, err := s.db.GetBookByID(ctx, bookID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &internal.ErrBookNotFound{}) {
			statusCode = http.StatusNotFound
		}

		return s.logAndReturnError(err, "failed to retrieve the book from the database", statusCode, logrus.Fields{"book_id": bookID})
	}

	// Set the book's new stsatus
	book.Status = newStatus

	// And save it
	updatedBook, err := s.db.UpdateBook(ctx, bookID, book)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &internal.ErrBookNotFound{}) {
			statusCode = http.StatusNotFound
		}

		return s.logAndReturnError(err, "failed to update the book in the database", statusCode, logrus.Fields{"book_id": bookID})
	}

	responseBody, err := json.Marshal(updatedBook)
	if err != nil {
		return s.logAndReturnError(err, "failed to encode the book into an http response", http.StatusInternalServerError, logrus.Fields{})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

func (s service) logAndReturnError(err error, message string, statusCode int, logFields logrus.Fields) (events.APIGatewayProxyResponse, error) {
	s.logger.WithError(err).WithFields(logFields).Error(message)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       formatErrorForResponseBody(fmt.Errorf("%s: %w", message, err)),
	}, nil
}

func formatErrorForResponseBody(err error) string {
	//TODO: Use a struct to represent the error
	//TODO: Allow this service to support Content-Type other than JSON
	return fmt.Sprintf(`{"error":"%s"}`, err.Error())
}

type errorResponse struct {
	ErrorMessage string `json:"error"`
}
