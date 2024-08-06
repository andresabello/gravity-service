# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the Go module files to download dependencies efficiently
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install Go tools
RUN go install github.com/air-verse/air@latest \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && go mod tidy

# Define the command to run your application with "Air" for live reloading
CMD ["air", "-c", "/app/.air.toml"]
