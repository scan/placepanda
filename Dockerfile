FROM golang:1.13-alpine AS builder

RUN apk add --update --upgrade ca-certificates build-base glib-dev vips-dev && \
    mkdir -p /src
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify

COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-w -s -h' -a -o /app

FROM alpine:3.10

RUN apk add --update --upgrade ca-certificates glib vips

COPY --from=builder /app /app

EXPOSE 8080

CMD ["/app"]
