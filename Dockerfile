FROM golang:1.25-alpine3.22 AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@v1.63.0

CMD ["/go/bin/air"]

FROM golang:1.25-alpine3.22 AS prd

ENV APP_PORT=3000

WORKDIR /app

COPY . .

RUN go mod download \
    && go mod verify \
    && go build -o main cmd/api/main.go \
    && go clean -cache -modcache -testcache

CMD ["/app/main"]

EXPOSE ${APP_PORT}
