FROM --platform=linux/amd64 golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./

# Download dependencies in advance; this will only re-run when the mod or sum files change.
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
RUN chown appuser:appgroup /app
USER appuser
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
