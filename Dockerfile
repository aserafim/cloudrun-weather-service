# Build stage
FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o cloudrun

# Final stage com Alpine para suporte TLS
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/cloudrun .

EXPOSE 8080

ENTRYPOINT ["./cloudrun"]