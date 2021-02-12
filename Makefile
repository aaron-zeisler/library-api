LAMBDA_SOURCE_DIR = ./lambdas
OUTPUT_DIR = dist
LAMBDA_OUTPUT_DIR := $(OUTPUT_DIR)/lambdas


.PHONY: test
test:
	@go test ./... -race -count=1 -cover


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


.PHONY: tidy
tidy:
	@go mod tidy


.PHONY: api
api: build
	@sam local start-api
