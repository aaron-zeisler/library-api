package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"

	"github.com/aaron-zeisler/library-api/internal"
	"github.com/aaron-zeisler/library-api/internal/testutils"
)

func Test_dynamodbBooksStorage_GetBooks(t *testing.T) {
	type state struct {
	}
	type expected struct {
		result []internal.Book
		err    error
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		/*
			"Testing Dynamo": {
				state{},
				expected{},
			},
		*/
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			s := NewDynamoDBBooksStorage()

			result, err := s.GetBooks(context.Background())
			fmt.Println(result)
			fmt.Println(err)

			assert.So(result, should.Resemble, tc.expected.result)
			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
		})
	}
}
