build:
	@go build -o ./bin/houndmail

run: build
	@./bin/houndmail

test:
	@go test -v ./...
