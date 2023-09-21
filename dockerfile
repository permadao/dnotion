FROM golang:1.18-alpine

WORKDIR /go/src/app

RUN apk add --no-cache bash

COPY . .

RUN go mod download
# backend
RUN go build -o dnotion ./start/main.go

EXPOSE 3002

# start up frontend
CMD ["./dnotion"]