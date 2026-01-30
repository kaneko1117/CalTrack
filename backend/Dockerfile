FROM golang:1.23-alpine

WORKDIR /app

# Install Air for hot-reload
RUN go install github.com/air-verse/air@v1.61.0

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Run with Air for hot-reload
CMD ["air", "-c", ".air.toml"]
