# Use an official Go runtime as a parent image
FROM golang:latest



RUN apt-get update; apt-get clean

# Install wget.
RUN apt-get install -y wget

RUN apt-get install -y gnupg

# Set the Chrome repo.
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list

# Install Chrome.
RUN apt-get update && apt-get -y install google-chrome-stable

# Set the working directory inside the container
WORKDIR /app

# Enable Go modules
ENV GO111MODULE=on

# Copy go.mod and go.sum to the container
COPY ./app/go.mod ./app/go.sum /app/

# Download dependencies
RUN go mod download

# Install "Air" for live reloading
RUN go install github.com/cosmtrek/air@latest

# Copy the rest of the application source code to the container
COPY ./app /app/

RUN go mod tidy

# Command to run the Go application using "Air" for live reloading
CMD ["air"]
