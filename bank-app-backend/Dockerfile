FROM golang:1.17-alpine AS builder

# Set the working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o drank-backend .

# Use a smaller image for the final stage
FROM alpine:3.15

# Install CA certificates for HTTPS connections and netcat for the startup script
RUN apk --no-cache add ca-certificates netcat-openbsd

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/drank-backend .
COPY --from=builder /app/.env .
COPY --from=builder /app/startup.sh .

# Make the startup script executable
RUN chmod +x startup.sh

# Expose the port the app runs on
EXPOSE 8080

# Command to run the startup script
CMD ["./startup.sh"]
