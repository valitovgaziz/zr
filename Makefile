build:
	@go build -o bin/zr cmd/main.go

run: build
	@./bin/zr

test:
	@go test ./...

.DEFAULT_GOAL=run