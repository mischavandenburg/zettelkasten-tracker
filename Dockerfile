# Use the official Golang image to create a build artifact.
FROM golang:1.22 as builder

# Set the working directory outside GOPATH to enable the support for modules.
WORKDIR /app

# Copy the go mod and sum files first to leverage Docker cache layering.
COPY go.mod go.sum ./
# Download dependencies in advance; this will only re-run when the mod or sum files change.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server .

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/server .

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable
CMD ["./server"]

