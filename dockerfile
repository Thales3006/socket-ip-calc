FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o server cmd/server/main.go

FROM debian:bookworm-slim

COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]