FROM golang:1.13 AS builder

RUN apt-get update -qq && \
    apt-get install -y ca-certificates build-essential libglib2.0-dev libvips-dev && \
    mkdir -p /src
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify

COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-w -s -h' -a -o /app

FROM debian:buster-slim

RUN apt-get update -qq && \
    apt-get install -y ca-certificates libglib2.0 libvips

COPY --from=builder /app /app
COPY ./pandas /pandas

EXPOSE 8080

CMD ["/app"]
