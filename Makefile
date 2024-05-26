build:
	@go build -o bin/redis-in-go .

run: build
	@./bin/redis-in-go
