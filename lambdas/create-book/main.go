package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"

	"github.com/aaron-zeisler/library-api/internal/books"
	"github.com/aaron-zeisler/library-api/internal/storage"
	"github.com/aaron-zeisler/library-api/lambdas"
)

func main() {
	db := storage.NewStaticBooksStorage()

	//TODO: Read these log settings from environment variables
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	service := books.NewService(db, books.WithLogger(logger))

	lambda.Start(lambdas.CORSWrapper(service.CreateBook))
}
