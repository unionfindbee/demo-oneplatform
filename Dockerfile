# Start from the Golang base image as a builder
FROM golang:1.21-rc-bookworm as builder

WORKDIR /build

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

COPY main.go .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app main.go

# Start from debian:bookworm-slim for the release image
FROM debian:bookworm-slim

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/app /app

# Expose port 7070 to the outside world
EXPOSE 7070

# Command to run the executable
CMD ["/app"]









# # Start from the Golang base image as a builder
# FROM golang:1.21-rc-bookworm as builder

# WORKDIR /

# COPY main2.go .

# # Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -cover -a -o app main2.go

# # Start from debian:bookworm-slim for the release image
# FROM debian:bookworm-slim

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app /app

# # Expose port 7070 to the outside world
# EXPOSE 7070

# # Command to run the executable
# CMD ["/app"]