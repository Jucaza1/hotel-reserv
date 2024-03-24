build:
	@go build -o bin/api main.go
run: build
	@./bin/api 
test:
	@go test -v ./...