package main

import (
	"github.com/aaron-zeisler/library-api/internal/books"
	"github.com/aaron-zeisler/library-api/internal/storage"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func main() {
	db := storage.NewStaticBooksStorage()
	logger := logrus.New()

	service := books.NewService(db, logger)

	lambda.Start(service.GetBooks)
}
