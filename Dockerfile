FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o /oauth-mock-server

# Make sure the binary is executable
RUN chmod +x /oauth-mock-server

# Command to run the executable
CMD ["/oauth-mock-server"]