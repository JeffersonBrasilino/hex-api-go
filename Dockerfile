FROM golang:1.24-alpine3.22 AS dev


RUN apk update && apk upgrade && apk add pkgconf \
git \
bash \
build-base \
sudo 

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest

RUN go mod download && go mod verify

CMD $(go env GOPATH)/bin/air

EXPOSE ${APP_PORT}


# FROM builder AS dev

# RUN go install github.com/air-verse/air@latest

# CMD $(go env GOPATH)/bin/air

# EXPOSE ${APP_PORT}

# FROM golang:1.23.4-alpine3.21 AS dev

# WORKDIR /app

# COPY . .

# CMD $(go env GOPATH)/bin/air
