build:
	@go build -o bin/redis-in-go .

run: build
	@./bin/redis-in-go --addr :5001

test:
	@go test -v ./...
