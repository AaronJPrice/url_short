local:
	@go run main/*.go

test:
	@go test -race ./...
