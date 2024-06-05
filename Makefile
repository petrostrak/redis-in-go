run: build
	@./bin/redis-in-go --listenAddr :5001
	
build:
	@go build -o bin/redis-in-go .

test:
	@go test -v ./...
