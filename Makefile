LAMBDA_SOURCE_DIR = ./lambdas
OUTPUT_DIR = dist
LAMBDA_OUTPUT_DIR := $(OUTPUT_DIR)/lambdas
AWS_TEMPLATE_FILE = template.yaml
AWS_PACKAGE_OUTPUT_FILE = packaged.yaml
AWS_REGION = us-west-1
S3_BUCKET = library-api-lambdas
CLOUDFORMATION_STACK_NAME = library-api-lambdas


.PHONY: test
test:
	@go test -race -count=1 -cover $$(go list ./... | grep -Ev 'mocks')


.PHONY: build
build: clean
	@go build ./...
	@for dir in `ls $(LAMBDA_SOURCE_DIR)`; do \
		GOOS=linux go build -o $(LAMBDA_OUTPUT_DIR)/$$dir $(LAMBDA_SOURCE_DIR)/$$dir; \
	done


.PHONY: clean
clean:
	@go clean ./...
	@rm -rf $(OUTPUT_DIR)
	@mkdir -p $(OUTPUT_DIR)
	@rm -f $(AWS_PACKAGE_OUTPUT_FILE)


.PHONY: tidy
tidy:
	@go mod tidy


.PHONY: mocks
mocks:
	@counterfeiter -o ./internal/books/mocks/mock_books_db.go --fake-name MockBooksDB ./internal/books booksDB


.PHONY: tools
tools:
	@go list -f '{{ join .Imports "\n" }}' -tags tools ./internal/tools | xargs go install


.PHONY: api
api: build
	@sam local start-api


.PHONY: package
package: build
	@sam package --template-file $(AWS_TEMPLATE_FILE) --s3-bucket $(S3_BUCKET) --region $(AWS_REGION) --output-template-file $(AWS_PACKAGE_OUTPUT_FILE)


.PHONY: deploy
deploy: package
	@sam deploy --template-file $(AWS_PACKAGE_OUTPUT_FILE) --stack-name $(CLOUDFORMATION_STACK_NAME) --capabilities CAPABILITY_IAM --region $(AWS_REGION)