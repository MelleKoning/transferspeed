# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the server application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server ./cmd/server

# Stage 2: Create the final, minimal image
FROM alpine:latest

COPY --from=builder /server /server

# Copy the data file
COPY data/lionface.jpg data/lionface.jpg
COPY data/lionheadbig.jpg data/lionheadbig.jpg

ENTRYPOINT ["/server"]
