# Stage 1: Build the application
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source file
COPY . /app/

# Compile the Go application
RUN go build -o logger .

# Stage 2: Setup the runtime container
FROM alpine

# Set the working directory in the new container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/logger /app/

# Copy the log files
COPY ./success.log /app/
COPY ./failure.log /app/

ENTRYPOINT ["./logger"]
