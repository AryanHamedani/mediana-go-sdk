FROM golang:1.20-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Disable IPv6 for go mod download to avoid "cannot assign requested address" errors
# See: https://groups.google.com/g/golang-nuts/c/KFZOOUeiYpc
RUN GODEBUG=netdns=cgo go mod download

# Copy source code
COPY . .

# Build the SDK
RUN go build -v ./...

# Set environment variables from .env file at runtime
CMD ["sh", "-c", "if [ -f .env ]; then export $(cat .env | grep -v '^#' | xargs); fi && cd examples && go run main.go"]

# Usage:
# Build: docker build -t mediana-sdk .
# Run: docker run --env-file .env mediana-sdk
# Or run with direct environment variables:
# docker run -e MEDIANA_API_KEY=your_key -e MEDIANA_TEST_PHONE=09XXXXXXXXX mediana-sdk 