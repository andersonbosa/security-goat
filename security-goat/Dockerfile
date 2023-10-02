# Use the official Go image as the base image
FROM golang:latest AS builder

LABEL org.opencontainers.image.source=https://github.com/andersonbosa/security-goat


# Set the working directory
WORKDIR /go/src/github.com/andersonbosa/security-goat

# Copy the source code into the container
COPY . .

# Build the Go application
RUN	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /security-goat main.go

# Use a minimal base image for the final image
FROM alpine:latest

# Copy the built binary from the builder image
COPY --from=builder /security-goat /security-goat

# Command to run the application with the provided arguments
ENTRYPOINT ["./security-goat"]

# Default arguments for the application (can be overridden when running the container)
# CMD ["--help"]
