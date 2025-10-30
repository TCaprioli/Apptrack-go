# Use the official Golang image as the base
FROM golang:1.23.0

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

# Make the startup script executable
COPY start.sh .
RUN chmod +x start.sh

# Expose port 8080
EXPOSE 8080

# Run the startup script
CMD ["./start.sh"]