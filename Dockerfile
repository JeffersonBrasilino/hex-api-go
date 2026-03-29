# 1. Dev Stage (Hot Reload)
FROM golang:1.25-alpine3.22 AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@v1.63.0
CMD ["/go/bin/air"]

# 2. Builder Stage (Compilation Edge)
FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app
# Copy dependencies specifier
COPY go.mod go.sum ./
RUN go mod download && go mod verify
# Copy application code
COPY . .
# Produce a statically-linked binary to strip away dependencies natively
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/api/main.go

# 3. Production Stage (Runtime Core)
FROM alpine:3.22 AS prd
# Install foundational operational dependencies (SSL and Timezones)
RUN apk --no-cache add ca-certificates tzdata
ENV APP_PORT=3000
WORKDIR /app
# Pull only the compiled artifact - drop the entire Golang ecosystem
COPY --from=builder /app/main .
RUN chmod +x ./main

EXPOSE ${APP_PORT}
CMD ["./main"]
