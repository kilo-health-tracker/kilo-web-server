FROM golang:1.19.0

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
CMD go run server.go -b 0.0.0.0