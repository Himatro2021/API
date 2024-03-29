FROM golang:1.17.1-alpine

WORKDIR /app

RUN mkdir bin
RUN mkdir src
RUN mkdir private_data

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod tidy

WORKDIR /app/src

COPY . .
RUN go build -o /app/bin/main main.go

WORKDIR /app/private_data

COPY private_data/ .

EXPOSE 8080

WORKDIR /app

RUN rm -r src/

CMD ["./bin/main", "server"]