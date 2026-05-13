# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway/main.go

# Final stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/gateway .
EXPOSE 8080
CMD ["./gateway"]
