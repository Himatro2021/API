FROM golang:1.17.1 as builder

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app


FROM alpine:3.14 as production

RUN apk add --no-cache ca-certificates

COPY --from=builder app .

EXPOSE 8080

CMD ./app