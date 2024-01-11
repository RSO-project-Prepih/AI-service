# Use official Golang image for building the binary
FROM golang:1.21.3 as builder

# Set working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-service .

# Use Alpine for the runtime image
FROM alpine:latest

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/ai-service /ai-service

# Expose the port that your app listens on
EXPOSE 8080

# Run the binary
CMD ["/ai-service"]