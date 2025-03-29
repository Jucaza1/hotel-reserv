build:
	@go build -o bin/api ./cmd/apiv1/main.go
run: build
	@./bin/api
seed:
	@go run ./seed/seed.go
seed-run: seed
	@./bin/api
test:
	@go test -v ./...

docker-mongo:
	@docker run --name mongodb -p 27017:27017 -d mongo:latest

docker-api:
	@echo "building Docker file"
	@docker build --no-cache -t api .
	@echo "running API inside Docker container"
	@docker run -p 4000:4000 api

docker-compose:
	@docker compose up --build -d
