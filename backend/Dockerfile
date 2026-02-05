FROM golang:1.24-alpine

WORKDIR /app

# Install Air for hot-reload
RUN go install github.com/air-verse/air@v1.61.0

# Install gotestsum for better test output
RUN go install gotest.tools/gotestsum@latest

# Install swag for Swagger doc generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Install sql-migrate for database migrations
RUN go install github.com/rubenv/sql-migrate/...@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Run with Air for hot-reload
CMD ["air", "-c", ".air.toml"]
