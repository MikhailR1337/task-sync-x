# Use an official Golang runtime as a parent image
FROM golang:1.20-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Build the Go app
RUN go build -o main ./cmd
# Expose port 3000 for the app to listen on
EXPOSE 3000

# Run the app when the container starts
CMD ["./main"]

