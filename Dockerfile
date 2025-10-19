# ------------------------------------------------------------
# Base build stage
# ------------------------------------------------------------

FROM golang:1.25.1-bookworm AS build_base

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached unless go.mod/go.sum changes)
RUN go mod download

# Code generation (if needed)
# RUN make generate

# ------------------------------------------------------------------------------
# App builder stage
# ------------------------------------------------------------------------------

FROM build_base AS app_builder


ADD . .

ENV GOPATH=/go
ARG GOARCH=arm64 # Default for Apple M class CPU

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go install ./main.go

# ------------------------------------------------------------
# Runner stage
# ------------------------------------------------------------

FROM alpine:3.19 AS app_runner

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy the pre-built binary file from the previous stage
COPY --from=app_builder /go/bin/main /usr/local/bin/fog


# Change ownership to non-root user
RUN chown appuser:appgroup /usr/local/bin/fog

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 7777

# Run the binary
CMD ["/usr/local/bin/fog"]