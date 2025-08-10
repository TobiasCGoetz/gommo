# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gommo .

# Final stage
FROM scratch

# Copy the binary from builder stage
COPY --from=builder /app/gommo /gommo

# Expose port (assuming the game runs on a port)
EXPOSE 8080

# Run the binary
CMD ["/gommo"]
