# Stage 1: Build the application
FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

# Use go mod download to download dependencies first
RUN go mod download

# Build the application
RUN go build -o main .

# Stage 2: Create a minimal image with only the compiled binary
FROM alpine

# Create a non-root user for running the application
RUN adduser -S -D -H -h appuser

USER appuser

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /build/main /app/

EXPOSE 3000

CMD ["./main"]