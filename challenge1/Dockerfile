# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/public_disclosure .

# ---- Runtime Stage ----
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary and runtime assets
COPY --from=builder /app/public_disclosure .
COPY --from=builder /build/conf ./conf
COPY --from=builder /build/static ./static
COPY --from=builder /build/views ./views

EXPOSE 3000

CMD ["./public_disclosure"]
