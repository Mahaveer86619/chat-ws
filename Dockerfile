# Use the official Go image as the base image
FROM golang:1.23-alpine as builder

# Set the current working directory in the container
WORKDIR /app

# Copy the Go modules files and download the dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o server .

# Create a new image for the final container
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory in the final container
WORKDIR /root/

# Copy the compiled Go binary from the builder image
COPY --from=builder /app/server .

# Expose the port the app will run on
EXPOSE 8080

# Run the compiled Go binary
CMD ["./server"]
