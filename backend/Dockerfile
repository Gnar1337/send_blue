# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install CA certs
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app .

# Final stage
FROM scratch
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]
EXPOSE 8080