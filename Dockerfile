FROM golang:1.17.1-alpine

WORKDIR /app

RUN mkdir bin
RUN mkdir src

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy

WORKDIR /app/src

COPY . .

RUN go build -o /app/bin/main main.go

WORKDIR /app

RUN rm -r src/

COPY .env .
COPY private_data/ .

EXPOSE 8080

CMD ["./bin/main", "server"]