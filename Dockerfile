FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/api ./cmd/apiv1/main.go
RUN go build -o bin/seed ./cmd/seed/seed.go

EXPOSE 4000

CMD ["/app/scripts/launch.sh"]
