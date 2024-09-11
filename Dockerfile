# Use the official Golang image as the base
FROM golang:1.22-alpine

# Set the working directory to root
WORKDIR /

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run the binary
CMD ["./main"]