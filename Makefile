build:
	@go build -o bin/api main.go
run: build
	@./bin/api 
seed:
	@go run scripts/seed.go
test:
	@go test -v ./...
docker:
	@docker run --name mongodb -p 27017:27017 -d mongo:latest