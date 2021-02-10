
test:
	@go test ./... -race -count=1 -cover

build:
	@go build ./...

clean:
	@go clean ./...

tidy:
	@go mod tidy

