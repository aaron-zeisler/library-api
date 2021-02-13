package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"

	"github.com/aaron-zeisler/library-api/internal"
	"github.com/aaron-zeisler/library-api/internal/testutils"
)

func TestNewStaticBookStorage(t *testing.T) {
	type state struct {
	}
	type expected struct {
		result staticBooksStorage
	}
	testCases := map[string]struct {
		state
		expected
	}{
		"Happy path is the only path": {
			state{},
			expected{
				result: staticBooksStorage{books: staticBooksData},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			result := NewStaticBooksStorage()

			assert.So(len(result.books), should.Equal, len(tc.expected.result.books))
			for id, book := range tc.expected.result.books {
				assert.So(result.books[id], should.Resemble, book)
			}
		})
	}
}

func Test_staticBookStorage_GetBooks(t *testing.T) {
	testBooks := map[string]internal.Book{
		"1": {ID: "1", Title: "Book 1"},
		"2": {ID: "2", Title: "Book 2"},
	}

	type state struct {
		books map[string]internal.Book
	}
	type expected struct {
		result []internal.Book
		err    error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"All book should be returned": {
			state{
				books: testBooks,
			},
			expected{
				result: []internal.Book{testBooks["1"], testBooks["2"]},
			},
		},
		"An empty library returns an empty slice": {
			state{
				books: map[string]internal.Book{},
			},
			expected{
				result: []internal.Book{},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			s := staticBooksStorage{
				books: tc.state.books,
			}

			result, err := s.GetBooks(context.Background())

			assert.So(len(result), should.Equal, len(tc.expected.result))
			for _, book := range tc.expected.result {
				assert.So(result, should.Contain, book)
			}
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_staticBookStorage_GetBookByID(t *testing.T) {
	testBooks := map[string]internal.Book{
		"1": {ID: "1", Title: "Book 1"},
		"2": {ID: "2", Title: "Book 2"},
	}

	type state struct {
		books  map[string]internal.Book
		bookID string
	}
	type expected struct {
		result internal.Book
		err    error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"Return the expected book": {
			state{
				books:  testBooks,
				bookID: "2",
			},
			expected{
				result: internal.Book{ID: "2", Title: "Book 2"},
			},
		},
		"An unknown book ID returns an error": {
			state{
				books:  testBooks,
				bookID: "7",
			},
			expected{
				result: internal.Book{},
				err:    internal.ErrBookNotFound{BookID: "7"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			s := staticBooksStorage{
				books: tc.state.books,
			}

			result, err := s.GetBookByID(context.Background(), tc.state.bookID)

			assert.So(result, should.Resemble, tc.expected.result)
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_staticBookStorage_CreateBook(t *testing.T) {
	type state struct {
		books       map[string]internal.Book // The libray's collection before the test
		title       string
		author      string
		isbn        string
		description string
	}
	type expected struct {
		result   internal.Book
		err      error
		numBooks int // The number of books that should be in the library after the test is run
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"Successfully create a new book": {
			state{
				books:       map[string]internal.Book{},
				title:       "Test book title",
				author:      "Test book author",
				isbn:        "Test book isbn",
				description: "Tests book description",
			},
			expected{
				result: internal.Book{
					Title:       "Test book title",
					Author:      "Test book author",
					ISBN:        "Test book isbn",
					Description: "Tests book description",
				},
				numBooks: 1,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			s := staticBooksStorage{
				books: tc.state.books,
			}

			result, err := s.CreateBook(context.Background(), tc.state.title, tc.state.author, tc.state.isbn, tc.state.description)

			// Verify the peropties of the Book object that was returned
			_, uuidErr := uuid.Parse(result.ID)
			assert.So(uuidErr, should.BeNil)

			assert.So(result.Title, should.Equal, tc.expected.result.Title)
			assert.So(result.Author, should.Equal, tc.expected.result.Author)
			assert.So(result.ISBN, should.Equal, tc.expected.result.ISBN)
			assert.So(result.Description, should.Equal, tc.expected.result.Description)

			// Verify that the book was added to the internal books collection
			assert.So(len(s.books), should.Equal, tc.expected.numBooks)
			if tc.expected.err == nil {
				assert.So(result, should.Resemble, s.books[result.ID])
			}

			// Verify the error if one was returned
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_staticBookStorage_UpdateBook(t *testing.T) {
	type state struct {
		books       map[string]internal.Book // The libray's collection before the test
		bookID      string
		title       string
		author      string
		isbn        string
		description string
	}
	type expected struct {
		result internal.Book
		err    error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"Successfully update a book's properties": {
			state{
				books: map[string]internal.Book{
					"448E55A3-E88E-4597-B3CB-11A844EFDA5D": {
						ID:          "448E55A3-E88E-4597-B3CB-11A844EFDA5D",
						Title:       "Fahrenheit 451",
						Author:      "Ray Bradbury",
						ISBN:        "9781451673265",
						Description: "It was a pleasure to burn",
					},
				},
				bookID:      "448E55A3-E88E-4597-B3CB-11A844EFDA5D",
				title:       "Something Else",
				author:      "Someone Else",
				isbn:        "9781451673265",
				description: "A different story altogether",
			},
			expected{
				result: internal.Book{
					ID:          "448E55A3-E88E-4597-B3CB-11A844EFDA5D",
					Title:       "Something Else",
					Author:      "Someone Else",
					ISBN:        "9781451673265",
					Description: "A different story altogether",
				},
			},
		},
		"An unknown book ID returns an error": {
			state{
				books:       make(map[string]internal.Book),
				bookID:      "448E55A3-E88E-4597-B3CB-11A844EFDA5D",
				title:       "Something Else",
				author:      "Someone Else",
				isbn:        "9781451673265",
				description: "A different story altogether",
			},
			expected{
				result: internal.Book{},
				err:    internal.ErrBookNotFound{BookID: "448E55A3-E88E-4597-B3CB-11A844EFDA5D"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			s := staticBooksStorage{
				books: tc.state.books,
			}

			result, err := s.UpdateBook(context.Background(), tc.state.bookID, tc.state.title, tc.state.author, tc.state.isbn, tc.state.description)

			// Verify the properties of the book object that was returned
			assert.So(result, should.Resemble, tc.expected.result)

			// Verify that the book is updated in the internal books collection
			if tc.expected.err == nil {
				assert.So(result, should.Resemble, s.books[result.ID])
			}

			// Verify the error if one was returned
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}

func Test_staticBookStorage_DeleteBook(t *testing.T) {
	type state struct {
		books  map[string]internal.Book // The libray's collection before the test
		bookID string
	}
	type expected struct {
		err      error
		numBooks int // The number of books that should be in the library after the test is run
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"Successfully delete a book": {
			state{
				books: map[string]internal.Book{
					"6B94AEF7-ABEC-483E-82CB-2B6E2B801997": {
						ID:          "6B94AEF7-ABEC-483E-82CB-2B6E2B801997",
						Title:       "Fear and Loathing in Las Vegas",
						Author:      "Hunter S. Thompson",
						ISBN:        "9780679785897",
						Description: "We were somewhere around Barstow on the edge of the desert when the drugs began to take hold",
					},
				},
				bookID: "6B94AEF7-ABEC-483E-82CB-2B6E2B801997",
			},
			expected{
				numBooks: 0,
			},
		},
		"An unknown book id returns an error": {
			state{
				books: map[string]internal.Book{
					"6B94AEF7-ABEC-483E-82CB-2B6E2B801997": {
						ID:          "6B94AEF7-ABEC-483E-82CB-2B6E2B801997",
						Title:       "Fear and Loathing in Las Vegas",
						Author:      "Hunter S. Thompson",
						ISBN:        "9780679785897",
						Description: "We were somewhere around Barstow on the edge of the desert when the drugs began to take hold",
					},
				},
				bookID: "11112222-3333-4444-5555-666677778888",
			},
			expected{
				err:      internal.ErrBookNotFound{BookID: "11112222-3333-4444-5555-666677778888"},
				numBooks: 1,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			s := staticBooksStorage{
				books: tc.state.books,
			}

			err := s.DeleteBook(context.Background(), tc.state.bookID)

			assert.So(err, testutils.ShouldEqualError, tc.expected.err)

			assert.So(len(s.books), should.Equal, tc.expected.numBooks)
		})
	}
}
