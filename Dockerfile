# Start from the Golang base image
FROM golang:1.20.5-bookworm

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Expose port 7070 to the outside world
EXPOSE 7070

# Command to run the executable
CMD ["./main"]
