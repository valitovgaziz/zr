build:
	@go build -o bin/zr cmd/main.go

run: build
	@./bin/zr

t:
	@go mod tidy

test:
	@go test ./...

.DEFAULT_GOAL=run