package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aaron-zeisler/library-api/internal"
)

type dynamodbBooksStorage struct {
	awsRegion string
	sess      *session.Session
	db        *dynamodb.DynamoDB
}

func NewDynamoDBBooksStorage(opts ...DynamoBooksStorageOption) *dynamodbBooksStorage {
	result := &dynamodbBooksStorage{
		awsRegion: "us-west-1", // Deftault region is us-west-1
	}

	for _, opt := range opts {
		opt(result)
	}

	awsConfig := &aws.Config{
		Region: aws.String(result.awsRegion),
	}

	result.sess = session.Must(session.NewSession(awsConfig))
	result.db = dynamodb.New(result.sess)

	return result
}

type DynamoBooksStorageOption func(*dynamodbBooksStorage)

func WithAWSRegion(awsRegion string) DynamoBooksStorageOption {
	return func(db *dynamodbBooksStorage) {
		db.awsRegion = awsRegion
	}
}

func (s *dynamodbBooksStorage) GetBooks(ctx context.Context) ([]internal.Book, error) {
	result := make([]internal.Book, 0)

	dbResult, err := s.db.ScanWithContext(ctx, &dynamodb.ScanInput{
		TableName: aws.String("library-api-books"),
	})
	if err != nil {
		return result, fmt.Errorf("failed to read from the databaseL %w", err)
	}

	err = dynamodbattribute.UnmarshalListOfMaps(dbResult.Items, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the result from the databaseL %w", err)
	}

	return result, nil
}

func (s *dynamodbBooksStorage) GetBookByID(ctx context.Context, bookID string) (internal.Book, error) {
	result := internal.Book{}
	return result, nil
}

func (s *dynamodbBooksStorage) CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error) {
	result := internal.Book{}
	return result, nil
}

func (s *dynamodbBooksStorage) UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error) {
	result := internal.Book{}
	return result, nil
}

func (s *dynamodbBooksStorage) DeleteBook(ctx context.Context, bookID string) error {
	return nil
}
