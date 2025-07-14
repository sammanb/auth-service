# Use minimal Go image
FROM golang:1.22-alpine

# Install git and curl for module fetching
RUN apk add --no-cache git curl

# Set working directory
WORKDIR /app

# Copy go mod and sun
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o auth-service main.go

# Expose PORT
EXPOSE 9000

# Start the app
CMD ["./auth-service"]