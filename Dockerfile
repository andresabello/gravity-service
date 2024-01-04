# Use an official Go runtime as a parent image
FROM golang:latest

RUN apt-get update; apt-get clean

RUN apt-get install -y wget

RUN apt-get install -y cron

# Set the working directory inside the container
WORKDIR /app

# Copy only the Go module files to download dependencies efficiently
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Enable Go modules
ENV GO111MODULE=on

# Install "Air" for live reloading
RUN go install github.com/cosmtrek/air@latest

RUN go mod tidy

# Set the working directory to the location of your application
WORKDIR /app/cmd/app

# Define the command to run your application with "Air" for live reloading
CMD ["air", "-c", "/app/.air.toml"]
