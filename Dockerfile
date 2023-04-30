FROM golang:1.19.0

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
RUN go build -o server.exe

CMD ./server.exe