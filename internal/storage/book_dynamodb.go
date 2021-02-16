package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"github.com/aaron-zeisler/library-api/internal"
)

type dynamodbBooksStorage struct {
	awsRegion string
	tableName string
	sess      *session.Session
	db        *dynamodb.DynamoDB
}

func NewDynamoDBBooksStorage(opts ...DynamoBooksStorageOption) *dynamodbBooksStorage {
	result := &dynamodbBooksStorage{
		awsRegion: "us-west-1", // Default region is us-west-1
		tableName: "library-api-books",
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
		return result, fmt.Errorf("failed to retrieve all the books from the database: %w", err)
	}

	err = dynamodbattribute.UnmarshalListOfMaps(dbResult.Items, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the result from the database: %w", err)
	}

	return result, nil
}

func (s *dynamodbBooksStorage) GetBookByID(ctx context.Context, bookID string) (internal.Book, error) {
	result := internal.Book{}

	key, err := dynamodbattribute.MarshalMap(map[string]string{"id": bookID})
	if err != nil {
		return result, fmt.Errorf("failed to marshal the bookID into a dynamo key: %w", err)
	}

	dbResult, err := s.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{TableName: aws.String(s.tableName), Key: key})
	if err != nil {
		return result, fmt.Errorf("failed to retrieve the book from the database: %w", err)
	}

	if len(dbResult.Item) == 0 {
		return result, internal.ErrBookNotFound{BookID: bookID}
	}

	err = dynamodbattribute.UnmarshalMap(dbResult.Item, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the result from the database: %w", err)
	}

	return result, nil
}

func (s *dynamodbBooksStorage) CreateBook(ctx context.Context, title, author, isbn, description string) (internal.Book, error) {
	result := internal.Book{}

	//TODO: Don't INSERT a book if the ISBN alreaedy exists in the database!
	newBook := internal.Book{
		ID:          uuid.New().String(),
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		Description: description,
	}
	item, err := dynamodbattribute.MarshalMap(newBook)
	if err != nil {
		return result, fmt.Errorf("failed to marshal the bookID into a dynamo key: %w", err)
	}

	_, err = s.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	})
	if err != nil {
		return result, fmt.Errorf("failed to create the new book in the database: %w", err)
	}

	return newBook, nil
}

func (s *dynamodbBooksStorage) UpdateBook(ctx context.Context, bookID, title, author, isbn, description string) (internal.Book, error) {
	result := internal.Book{}

	// Attempt to retrieve the book to verify that it exists before updating
	_, err := s.GetBookByID(ctx, bookID)
	if err != nil {
		return result, err
	}

	key, err := dynamodbattribute.MarshalMap(map[string]string{"id": bookID})
	if err != nil {
		return result, fmt.Errorf("failed to marshal the bookID into a dynamo key: %w", err)
	}

	book := struct {
		Title       string `json:":t"`
		Author      string `json:":a"`
		ISBN        string `json:":i"`
		Description string `json:":d"`
	}{
		Title:       title,
		Author:      author,
		ISBN:        isbn,
		Description: description,
	}
	updates, err := dynamodbattribute.MarshalMap(book)
	if err != nil {
		return result, fmt.Errorf("failed to marshal the bookID into a dynamo key: %w", err)
	}

	dbResult, err := s.db.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.tableName),
		Key:                       key,
		UpdateExpression:          aws.String("SET isbn=:i, title=:t, author=:a, description=:d"),
		ExpressionAttributeValues: updates,
		ReturnValues:              aws.String("ALL_NEW"),
	})
	if err != nil {
		return result, fmt.Errorf("failed to update the book in the database: %w", err)
	}

	if len(dbResult.Attributes) == 0 {
		return result, internal.ErrBookNotFound{BookID: bookID}
	}

	err = dynamodbattribute.UnmarshalMap(dbResult.Attributes, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the result from the database: %w", err)
	}

	return result, nil
}

func (s *dynamodbBooksStorage) DeleteBook(ctx context.Context, bookID string) error {
	// Attempt to retrieve the book to verify that it exists before deleting
	_, err := s.GetBookByID(ctx, bookID)
	if err != nil {
		return err
	}

	key, err := dynamodbattribute.MarshalMap(map[string]string{"id": bookID})
	if err != nil {
		return fmt.Errorf("failed to marshal the bookID into a dynamo key: %w", err)
	}

	_, err = s.db.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete the book from the database: %w", err)
	}

	return nil
}
