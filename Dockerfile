# --- Stage 1: Build Rust Core ---
FROM rust:1.84-slim-bookworm AS rust-builder
WORKDIR /app/core
COPY core/ .
RUN cargo build --release

# --- Stage 2: Build Go CLI ---
FROM golang:1.24-bookworm AS go-builder
WORKDIR /app/cli
COPY cli/ .
# Note: We don't embed the binary yet, we just build the CLI
RUN go build -o shadowprism main.go

# --- Stage 3: Final Image ---
FROM debian:bookworm-slim
WORKDIR /app

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Copy binaries from builders
COPY --from=rust-builder /app/core/target/release/shadowprism-core /app/shadowprism-core
COPY --from=go-builder /app/cli/shadowprism /app/shadowprism

# Create the data directory
RUN mkdir -p /root/.shadowprism

# Set environment variables
ENV SHADOWPRISM_AUTH_TOKEN=container-token-xyz
ENV PATH="/app:${PATH}"

# The container will run the Go CLI by default
# For Docker, we use the bot mode or a long-running TUI mode
ENTRYPOINT ["shadowprism"]
CMD ["bot"]
