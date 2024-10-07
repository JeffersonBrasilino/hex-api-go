FROM golang:latest AS dev

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

EXPOSE ${APP_PORT}

CMD $(go env GOPATH)/bin/air
