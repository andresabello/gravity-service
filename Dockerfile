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

# Install Chromium (open-source version of Chrome)
RUN apt-get update && \
    apt-get install -y \
    chromium \
    libnss3 \
    libglib2.0-0 \
    libx11-6 \
    libx11-xcb1 \
    libxcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    ca-certificates \
    fonts-liberation \
    libappindicator1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcairo2 \
    libcups2 \
    libdbus-1-3 \
    libfontconfig1 \
    libgbm1 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libxcomposite1 \
    libxdamage1 \
    libxkbcommon0 \
    libxslt1.1 \
    libxss1 \
    libxtst6 \
    xdg-utils

# Set the Chrome path
ENV CHROME_BIN=/usr/bin/chromium-browser

# Install "Air" for live reloading
RUN go install github.com/air-verse/air@latest

# Install migrate package. Use postgress
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install chromedp
RUN go get -u github.com/chromedp/chromedp

RUN go mod tidy

# Set the working directory to the location of your application
WORKDIR /app

# Define the command to run your application with "Air" for live reloading
CMD ["air", "-c", "/app/.air.toml"]
