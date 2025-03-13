# Use Go release candidate image since dependencies require Go 1.23+
FROM golang:rc-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o main .

# Create a non-root user to run the application
RUN adduser -D -g '' appuser
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 