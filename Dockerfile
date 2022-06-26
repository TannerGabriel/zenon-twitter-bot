FROM golang:1.17

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apt-get update && apt-get upgrade -y
RUN apt-get install libzmq3-dev -y
RUN apt-get install libczmq-dev -y

RUN export GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

# Build the application
RUN go build -o main cmd/main.go

# Command to run the executable
CMD ["./main"]