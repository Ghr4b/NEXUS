# Build Stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /build
COPY . .
RUN go build -o /app-binary .

# Final Stage
FROM alpine:3.19

# Create a non-root user and group
RUN addgroup -S ctfgroup && adduser -S ctfuser -G ctfgroup

# Install curl for healthcheck
RUN apk --no-cache add curl

# Create required directories
RUN mkdir /app /logs

# Copy artifacts from builder
COPY --from=builder /app-binary /app/server
COPY flag.txt /flag.txt
COPY logs/ /logs/

# Ensure strict permissions
RUN chown -R root:ctfgroup /app /logs /flag.txt && \
  chmod 550 /app/server && \
  chmod 440 /flag.txt && \
  find /logs -type d -exec chmod 550 {} + && \
  find /logs -type f -exec chmod 440 {} +

# Switch to non-root user
USER ctfuser

# Expose the correct port
EXPOSE 8080

# Configure environment variables
ENV logDir=/logs


HEALTHCHECK --interval=1m --timeout=10s --retries=3 \
  CMD curl -sf http://localhost:9091/health | grep -qi '"status":"UP"' || exit 1


CMD ["/app/server"]
