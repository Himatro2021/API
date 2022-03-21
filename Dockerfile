FROM golang:1.17.1-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy

COPY . .

RUN mkdir bin

RUN go build -o ./bin/main main.go

EXPOSE 8080

CMD ["./bin/main", "--add-host=host.docker.internal:host-gateway"]