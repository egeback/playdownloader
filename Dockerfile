# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Marky Egebäck <marky@egeback.se>"

RUN apt-get update && apt-get install -y python3 python3-pip
RUN apt-get install -y svtplay-dl

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
#COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# RUN go build -o main internal/main.go
#RUN ./cmd/build.sh

# Expose port 8080 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
