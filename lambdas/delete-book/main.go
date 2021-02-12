package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"

	"github.com/aaron-zeisler/library-api/internal/books"
	"github.com/aaron-zeisler/library-api/internal/storage"
)

func main() {
	db := storage.NewStaticBooksStorage()
	logger := logrus.New()

	service := books.NewService(db, books.WithLogger(logger))

	lambda.Start(service.DeleteBook)
}
