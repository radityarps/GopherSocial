FROM golang:alpine

# Install git and curl
RUN apk add --no-cache git curl && \
    git config --global --add safe.directory /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Expose default HTTP port
EXPOSE 8080

# Run Air live reload
CMD ["air", "-c", ".air.toml"]
