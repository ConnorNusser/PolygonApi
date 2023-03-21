build:
	@go build -o bin/polygonApi
run: build
	@./bin/polygonApi
test:
	@go test -v ./...