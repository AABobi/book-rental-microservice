# Use the official Golang image as the base image
FROM golang:1.21.3

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o /book ./cmd/api

# Expose the port that the application will run on
EXPOSE 80

# Define the command to run your application
CMD ["/book"]