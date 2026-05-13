#!/bin/bash
set -e

# Project configuration
PROJECT_DIR="$HOME/go-assistant"

# Create the directory structure
mkdir -p "$PROJECT_DIR/cmd/gateway"
mkdir -p "$PROJECT_DIR/internal/gateway"
mkdir -p "$PROJECT_DIR/internal/tasks"
mkdir -p "$PROJECT_DIR/internal/location"
mkdir -p "$PROJECT_DIR/internal/notifications"
mkdir -p "$PROJECT_DIR/internal/integrations"
mkdir -p "$PROJECT_DIR/internal/ai"
mkdir -p "$PROJECT_DIR/pkg/models"
mkdir -p "$PROJECT_DIR/deploy/docker"
mkdir -p "$PROJECT_DIR/deploy/k8s"
mkdir -p "$PROJECT_DIR/deploy/terraform"

cd "$PROJECT_DIR"

# Write go.mod
cat > go.mod << 'EOF'
module github.com/go-assistant/core

go 1.22
EOF

# Write cmd/gateway/main.go
cat > cmd/gateway/main.go << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Task represents a task item
type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var tasks = []Task{
	{ID: "1", Title: "Initialize project structure"},
	{ID: "2", Title: "Setup Docker and Compose"},
}

func main() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// List tasks endpoint
	mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	// Create task endpoint
	mux.HandleFunc("POST /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if task.ID == "" || task.Title == "" {
			http.Error(w, "ID and Title are required", http.StatusBadRequest)
			return
		}
		tasks = append(tasks, task)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
	})

	const port = ":8080"
	log.Printf("Gateway server starting on %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
EOF

# Write docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_URL=postgres://goassistant:secret@postgres:5432/goassistant?sslmode=disable
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: goassistant
      POSTGRES_USER: goassistant
      POSTGRES_PASSWORD: secret
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
EOF

# Write Dockerfile
cat > Dockerfile << 'EOF'
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
EOF

# Write Makefile
cat > Makefile << 'EOF'
.PHONY: build run docker-up docker-down test

build:
	go build -o bin/gateway ./cmd/gateway/main.go

run:
	go run ./cmd/gateway/main.go

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

test:
	go test ./...
EOF

# Write .env.example
cat > .env.example << 'EOF'
NOTION_API_KEY=your_notion_key
GMAIL_CLIENT_ID=your_gmail_id
GMAIL_CLIENT_SECRET=your_gmail_secret
GOOGLE_CALENDAR_ID=your_calendar_id
OPENAI_API_KEY=your_openai_key
POSTGRES_URL=postgres://goassistant:secret@localhost:5432/goassistant?sslmode=disable
REDIS_URL=redis://localhost:6379
EOF

echo "=== go-assistant initialized ==="
echo "Next steps:"
echo "1. cd ~/go-assistant"
echo "2. cp .env.example .env"
echo "3. make docker-up"
echo "4. Verify health: curl http://localhost:8080/health"
