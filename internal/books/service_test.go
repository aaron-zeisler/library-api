package books

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aaron-zeisler/library-api/internal"
	"github.com/aaron-zeisler/library-api/internal/books/mocks"
	"github.com/aaron-zeisler/library-api/internal/testutils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func Test_service_GetBooks(t *testing.T) {
	type state struct {
		request    events.APIGatewayProxyRequest
		dbResponse []internal.Book
		dbError    error
	}
	type expected struct {
		responseCode int
		responseBody interface{}
		err          error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"The call to db.GetBooks returns an error": {
			state{
				request:    events.APIGatewayProxyRequest{},
				dbResponse: nil,
				dbError:    errors.New("db.GetBooks error"),
			},
			expected{
				responseCode: http.StatusInternalServerError,
				responseBody: errorResponse{
					ErrorMessage: "failed to retrieve books from the database: db.GetBooks error",
				},
			},
		},
		"Happy path": {
			state{
				request: events.APIGatewayProxyRequest{},
				dbResponse: []internal.Book{
					{ID: "12345", ISBN: "12345", Title: "GetBooks Test"},
				},
			},
			expected{
				responseCode: http.StatusOK,
				responseBody: []internal.Book{
					{ID: "12345", ISBN: "12345", Title: "GetBooks Test"},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			db := &mocks.MockBooksDB{}
			db.GetBooksReturns(tc.state.dbResponse, tc.state.dbError)

			s := service{
				db:     db,
				logger: logrus.New(),
			}

			result, err := s.GetBooks(context.Background(), tc.state.request)

			// Verify the response code
			assert.So(result.StatusCode, should.Equal, tc.expected.responseCode)

			// Verify the response body
			if tc.expected.responseCode == http.StatusOK {
				resp := []internal.Book{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			} else {
				resp := errorResponse{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			}

			// Verify the error
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_service_GetBookByID(t *testing.T) {
	type state struct {
		request    events.APIGatewayProxyRequest
		dbResponse internal.Book
		dbError    error
	}
	type expected struct {
		responseCode int
		responseBody interface{}
		err          error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"db.GetBookByID returns a BookNotFound error": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
				},
				dbResponse: internal.Book{},
				dbError:    internal.ErrBookNotFound{BookID: "12345"},
			},
			expected{
				responseCode: http.StatusNotFound,
				responseBody: errorResponse{
					ErrorMessage: "failed to retrieve the book from the database: The book with ID '12345' was not found",
				},
			},
		},
		"db.GetBookByID returns an unepected error": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
				},
				dbResponse: internal.Book{},
				dbError:    errors.New("db.GetBookByID error"),
			},
			expected{
				responseCode: http.StatusInternalServerError,
				responseBody: errorResponse{
					ErrorMessage: "failed to retrieve the book from the database: db.GetBookByID error",
				},
			},
		},
		"Happy path": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
				},
				dbResponse: internal.Book{
					ID: "12345", ISBN: "12345", Title: "GetBookByID Test",
				},
			},
			expected{
				responseCode: http.StatusOK,
				responseBody: internal.Book{
					ID: "12345", ISBN: "12345", Title: "GetBookByID Test",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			db := &mocks.MockBooksDB{}
			db.GetBookByIDReturns(tc.state.dbResponse, tc.state.dbError)

			s := service{
				db:     db,
				logger: logrus.New(),
			}

			result, err := s.GetBookByID(context.Background(), tc.state.request)

			// Verify the response code
			assert.So(result.StatusCode, should.Equal, tc.expected.responseCode)

			// Verify the response body
			if tc.expected.responseCode == http.StatusOK {
				resp := internal.Book{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			} else {
				resp := errorResponse{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			}

			// Verify the error
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_service_CreateBook(t *testing.T) {
	type state struct {
		request    events.APIGatewayProxyRequest
		dbResponse internal.Book
		dbError    error
	}
	type expected struct {
		responseCode int
		responseBody interface{}
		err          error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"The request body is malformed": {
			state{
				request: events.APIGatewayProxyRequest{
					Body: `}`,
				},
			},
			expected{
				responseCode: http.StatusBadRequest,
				responseBody: errorResponse{
					ErrorMessage: "failed to decode the request body into a book object: invalid character '}' looking for beginning of value",
				},
			},
		},
		"db.CreateBook returns an error": {
			state{
				request: events.APIGatewayProxyRequest{
					Body: `{"isbn": "12345", "title": "CreateBook Test", "author": "Testy McTesterson"}`,
				},
				dbResponse: internal.Book{},
				dbError:    errors.New("db.CreateBook error"),
			},
			expected{
				responseCode: http.StatusInternalServerError,
				responseBody: errorResponse{
					ErrorMessage: "failed to create a new book in the database: db.CreateBook error",
				},
			},
		},
		"Happy path": {
			state{
				request: events.APIGatewayProxyRequest{
					Body: `{"isbn": "12345", "title": "CreateBook Test", "author": "Testy McTesterson"}`,
				},
				dbResponse: internal.Book{
					ID: "12345", ISBN: "12345", Title: "CreateBook Title", Author: "Testy McTesterson",
				},
			},
			expected{
				responseCode: http.StatusOK,
				responseBody: internal.Book{
					ID: "12345", ISBN: "12345", Title: "CreateBook Title", Author: "Testy McTesterson",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			db := &mocks.MockBooksDB{}
			db.CreateBookReturns(tc.state.dbResponse, tc.state.dbError)

			s := service{
				db:     db,
				logger: logrus.New(),
			}

			result, err := s.CreateBook(context.Background(), tc.state.request)

			// Verify the response code
			assert.So(result.StatusCode, should.Equal, tc.expected.responseCode)

			// Verify the response body
			if tc.expected.responseCode == http.StatusOK {
				resp := internal.Book{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			} else {
				resp := errorResponse{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			}

			// Verify the error
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_service_UpdateBook(t *testing.T) {
	type state struct {
		request    events.APIGatewayProxyRequest
		dbResponse internal.Book
		dbError    error
	}
	type expected struct {
		responseCode int
		responseBody interface{}
		err          error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"The request body is malformed": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
					Body:           `}`,
				},
			},
			expected{
				responseCode: http.StatusBadRequest,
				responseBody: errorResponse{
					ErrorMessage: "failed to decode the request body into a book object: invalid character '}' looking for beginning of value",
				},
			},
		},
		"db.UpdateBook returns a BookNotFound error": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
					Body:           `{"id": "12345", "isbn": "12345", "title": "UpdateBook Test", "author": "Testy McTesterson"}`,
				},
				dbResponse: internal.Book{},
				dbError:    internal.ErrBookNotFound{BookID: "12345"},
			},
			expected{
				responseCode: http.StatusNotFound,
				responseBody: errorResponse{
					ErrorMessage: "failed to update the book in the database: The book with ID '12345' was not found",
				},
			},
		},
		"db.UpdateBook returns an unexpected error": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
					Body:           `{"id": "12345", "isbn": "12345", "title": "UpdateBook Test", "author": "Testy McTesterson"}`,
				},
				dbResponse: internal.Book{},
				dbError:    errors.New("db.UpdateBook error"),
			},
			expected{
				responseCode: http.StatusInternalServerError,
				responseBody: errorResponse{
					ErrorMessage: "failed to update the book in the database: db.UpdateBook error",
				},
			},
		},
		"Happy path": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
					Body:           `{"id": "12345", "isbn": "12345", "title": "UpdateBook Test", "author": "Testy McTesterson"}`,
				},
				dbResponse: internal.Book{
					ID: "12345", ISBN: "12345", Title: "UpdateBook Test", Author: "Testy McTesterson",
				},
			},
			expected{
				responseCode: http.StatusOK,
				responseBody: internal.Book{
					ID: "12345", ISBN: "12345", Title: "UpdateBook Test", Author: "Testy McTesterson",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			db := &mocks.MockBooksDB{}
			db.UpdateBookReturns(tc.state.dbResponse, tc.state.dbError)

			s := service{
				db:     db,
				logger: logrus.New(),
			}

			result, err := s.UpdateBook(context.Background(), tc.state.request)

			// Verify the response code
			assert.So(result.StatusCode, should.Equal, tc.expected.responseCode)

			// Verify the response body
			if tc.expected.responseCode == http.StatusOK {
				resp := internal.Book{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			} else {
				resp := errorResponse{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			}

			// Verify the error
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_service_DeleteBook(t *testing.T) {
	type state struct {
		request events.APIGatewayProxyRequest
		dbError error
	}
	type expected struct {
		responseCode int
		responseBody interface{}
		err          error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"db.DeleteBook returns a BookNotFound error": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
				},
				dbError: internal.ErrBookNotFound{BookID: "12345"},
			},
			expected{
				responseCode: http.StatusOK,
			},
		},
		"db.DeleteBook returns an unexpected error": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
				},
				dbError: errors.New("db.DeleteBook error"),
			},
			expected{
				responseCode: http.StatusInternalServerError,
				responseBody: errorResponse{
					ErrorMessage: "failed to delete the book from the database: db.DeleteBook error",
				},
			},
		},
		"Happy path": {
			state{
				request: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"book_id": "12345"},
				},
			},
			expected{
				responseCode: http.StatusOK,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			db := &mocks.MockBooksDB{}
			db.DeleteBookReturns(tc.state.dbError)

			s := service{
				db:     db,
				logger: logrus.New(),
			}

			result, err := s.DeleteBook(context.Background(), tc.state.request)

			// Verify the response code
			assert.So(result.StatusCode, should.Equal, tc.expected.responseCode)

			// Verify the response body
			if tc.expected.responseCode != http.StatusOK {
				resp := errorResponse{}
				jsonErr := json.Unmarshal([]byte(result.Body), &resp)
				assert.So(jsonErr, should.BeNil)
				assert.So(resp, should.Resemble, tc.expected.responseBody)
			}

			// Verify the error
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}
