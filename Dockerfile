# Start from the latest golang base image
FROM golang:1.23-alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
WORKDIR /app/cmd/nats-readiness
RUN go build -o /nats-readiness

FROM alpine:3.20

WORKDIR /app/

COPY --from=builder /nats-readiness .

ENTRYPOINT ["./nats-readiness"]
