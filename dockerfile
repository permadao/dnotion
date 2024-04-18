FROM golang:1.21-alpine

WORKDIR /go/src/app

RUN apk add --no-cache bash

COPY . .

RUN go mod download
# backend
RUN go build -o dnotion ./run/service

# Copy config file
COPY config/config.toml .

# start up frontend
CMD ["./dnotion"]