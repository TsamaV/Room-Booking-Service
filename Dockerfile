FROM golang:1.24.4-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]