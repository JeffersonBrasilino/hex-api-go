FROM golang:1.24-alpine3.22 AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

CMD ["/go/bin/air"]

FROM golang:1.24-alpine3.22 AS prd

ENV APP_PORT=3000

WORKDIR /app

COPY . .

RUN go mod download \
    && go mod verify \
    && go build -o main cmd/api/main.go \
    && go clean -cache -modcache -testcache

CMD ["/app/main"]

EXPOSE ${APP_PORT}
