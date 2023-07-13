# Start from the Golang base image as a builder
FROM golang:1.13.3-buster as builder

WORKDIR /

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the builder container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app main.go

# # Start from debian:bookworm-slim for the release image
# FROM debian:bookworm-slim

# # Copy the binary from the builder container to the release container
# COPY --from=builder /app/ /app/

# # Expose port 7070 to the outside world
# EXPOSE 7070

# Command to run the executable
CMD ["/app"]