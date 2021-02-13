package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/aaron-zeisler/library-api/internal"
)

type staticBooksStorage struct {
	books map[string]internal.Book
}

func NewStaticBooksStorage() *staticBooksStorage {
	return &staticBooksStorage{
		books: staticBooksData,
	}
}

func (s *staticBooksStorage) GetBooks(ctx context.Context) ([]internal.Book, error) {
	result := make([]internal.Book, 0, len(s.books))
	for _, book := range s.books {
		result = append(result, book)
	}
	return result, nil
}

func (s *staticBooksStorage) GetBookByID(ctx context.Context, bookID string) (internal.Book, error) {
	book, ok := s.books[bookID]
	if !ok {
		return internal.Book{}, internal.ErrBookNotFound{BookID: bookID}
	}
	return book, nil
}

func (s *staticBooksStorage) CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error) {
	newBookID := uuid.New().String()
	newBook := internal.Book{
		ID:          newBookID,
		Title:       title,
		Author:      author,
		ISBN:        isbn,
		Description: description,
	}

	s.books[newBookID] = newBook

	return newBook, nil
}

func (s *staticBooksStorage) UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error) {
	_, ok := s.books[bookID]
	if !ok {
		return internal.Book{}, internal.ErrBookNotFound{BookID: bookID}
	}

	s.books[bookID] = internal.Book{
		ID:          bookID,
		Title:       title,
		Author:      author,
		ISBN:        isbn,
		Description: description,
	}
	return s.books[bookID], nil

}

func (s *staticBooksStorage) DeleteBook(ctx context.Context, bookID string) error {
	_, ok := s.books[bookID]
	if !ok {
		return internal.ErrBookNotFound{BookID: bookID}
	}

	delete(s.books, bookID)

	return nil
}

var staticBooksData = map[string]internal.Book{
	"448E55A3-E88E-4597-B3CB-11A844EFDA5D": {
		ID:          "448E55A3-E88E-4597-B3CB-11A844EFDA5D",
		Title:       "Fahrenheit 451",
		Author:      "Ray Bradbury",
		ISBN:        "9781451673265",
		Description: "It was a pleasure to burn",
	},
	"2E09FDF3-9DF2-4320-A86E-E2178262D4E6": {
		ID:          "2E09FDF3-9DF2-4320-A86E-E2178262D4E6",
		Title:       "1984",
		Author:      "George Orwell",
		ISBN:        "9780452284234",
		Description: "It was a bright cold day in April, and the clocks were striking thirteen",
	},
	"9AEFE32B-0B69-4D9D-BC7B-E4B1C2B8616D": {
		ID:          "9AEFE32B-0B69-4D9D-BC7B-E4B1C2B8616D",
		Title:       "Anna Karenina",
		Author:      "Leo Tolstoy",
		ISBN:        "9798560833640",
		Description: "Happy families are all alike; every unhappy family is unhappy in its own way",
	},
	"EA31C594-CDBB-4740-981F-B77AFA3C0FBA": {
		ID:          "EA31C594-CDBB-4740-981F-B77AFA3C0FBA",
		Title:       "Moby Dick",
		Author:      "Herman Melville",
		ISBN:        "9781514649749",
		Description: "Call me Ishmael",
	},
	"E7EC1121-8310-4B1D-93C2-BFF01B5F90A2": {
		ID:          "E7EC1121-8310-4B1D-93C2-BFF01B5F90A2",
		Title:       "The Great Gatsby",
		Author:      "F. Scott Fitzgerald",
		ISBN:        "9780743273565",
		Description: "In my younger an more vulnerable years my father gave me some advice that I've been turning over in my mind ever since",
	},
	"E0805B34-2369-469F-9AE0-81A812229A86": {
		ID:          "E0805B34-2369-469F-9AE0-81A812229A86",
		Title:       "The Catcher in the Rye",
		Author:      "J.D. Salinger",
		ISBN:        "9780316769174",
		Description: "If you really want to hear about it, the first thing you’ll probably want to know is where I was born, and what my lousy childhood was like, and how my parents were occupied and all before they had me, and all that David Copperfield kind of crap, but I don’t feel like going into it, if you want to know the truth",
	},
	"0E119988-56A7-487B-AC3A-C867CC4D4353": {
		ID:          "0E119988-56A7-487B-AC3A-C867CC4D4353",
		Title:       "The Restaurant at the End of the Universe",
		Author:      "Douglas Adams",
		ISBN:        "9789123918430",
		Description: "The story so far: in the beginning, the universe was created. This has made a lot of people very angry and been widely regarded as a bad move",
	},
	"6B94AEF7-ABEC-483E-82CB-2B6E2B801997": {
		ID:          "6B94AEF7-ABEC-483E-82CB-2B6E2B801997",
		Title:       "Fear and Loathing in Las Vegas",
		Author:      "Hunter S. Thompson",
		ISBN:        "9780679785897",
		Description: "We were somewhere around Barstow on the edge of the desert when the drugs began to take hold",
	},
	"3E020259-42AF-4564-BF1F-FC57B0977EE2": {
		ID:          "3E020259-42AF-4564-BF1F-FC57B0977EE2",
		Title:       "Beloved",
		Author:      "Toni Morrison",
		ISBN:        "9781400033416",
		Description: "124 was spiteful. Full of Baby's venom",
	},
	"78D9D95E-EE03-4B0D-8DAA-5E0BB9AC11D7": {
		ID:          "78D9D95E-EE03-4B0D-8DAA-5E0BB9AC11D7",
		Title:       "The Martian",
		Author:      "Andy Weir",
		ISBN:        "9781101905005",
		Description: "I'm pretty much f*cked",
	},
	"33BBCE73-BABF-40D1-BFF9-55DEC242BBEE": {
		ID:          "33BBCE73-BABF-40D1-BFF9-55DEC242BBEE",
		Title:       "Harry Potter and the Sorceror's Stone",
		Author:      "J.K. Rowling",
		ISBN:        "9781338596700",
		Description: "Mr and Mrs Dursley, of number four Privet Drive, were proud to say that they were perfectly normal, thank you very much",
	},
	"B98C89F1-E8F6-43FD-A3F8-4D5A1DA8E30B": {
		ID:          "B98C89F1-E8F6-43FD-A3F8-4D5A1DA8E30B",
		Title:       "Old Man's War",
		Author:      "John Scalzi",
		ISBN:        "9780765348272",
		Description: "On his 75th birthday John Perry did two things. First, he visited his wife’s grave. Then he joined the army",
	},
}
