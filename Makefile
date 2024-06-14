build:
	@go build -o bin/api main.go
run: build
	@./bin/api 
seed:
	@go run scripts/seed.go
test:
	@go test -v ./...

mongo:
	@docker run --name mongodb -p 27017:27017 -d mongo:latest
docker:
	echo "building Docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 4000:4000 api
docker-compose:
	@docker-compose up -d