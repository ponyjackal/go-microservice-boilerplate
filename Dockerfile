# Start from golang base image
FROM golang:1.21-alpine as builder

# Install git and vips dependencies
RUN apk update && apk add --no-cache git alpine-sdk vips-dev

# Working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy everything
COPY . .

# Build the Go app
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=1 GOOS=linux go build -mod=readonly -v -o main .

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates and vips runtime libraries
RUN apk --no-cache add ca-certificates vips

WORKDIR /root/

# Copy the pre-built binary file from the previous stage. Also copy config yml file
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8000
EXPOSE 50051

# Command to run the executable
CMD ["./main"]
