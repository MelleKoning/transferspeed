# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the client application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /client ./cmd/client

# Stage 2: Create the final, minimal image
FROM alpine:latest

COPY --from=builder /client /client

ENTRYPOINT ["/client"]
