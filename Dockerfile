# TinyClaw (tinyland-inc/tinyclaw) — standalone Dockerfile
#
# Builds the TinyClaw-based agent with RemoteJuggler config:
# - Multi-stage Go build
# - Bakes in a config.json with Aperture API routing
# - Health check on /health port 18790
#
# Build context: repo root
# GHCR workflow builds from main branch pushes.

# ============================================================
# Stage 1: Build the tinyclaw binary
# ============================================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 go generate ./... && \
    CGO_ENABLED=0 go build -v -tags stdjson \
      -ldflags "-X github.com/tinyland-inc/tinyclaw/cmd/tinyclaw/internal.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev) -s -w" \
      -o build/tinyclaw ./cmd/tinyclaw

# ============================================================
# Stage 2: Minimal runtime image
# ============================================================
FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata curl

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q --spider http://localhost:18790/health || exit 1

COPY --from=builder /src/build/tinyclaw /usr/local/bin/tinyclaw

RUN addgroup -g 1000 tinyclaw && \
    adduser -D -u 1000 -G tinyclaw tinyclaw && \
    mkdir -p /workspace && chown tinyclaw:tinyclaw /workspace

USER tinyclaw

# Run onboard to create initial directories and config
RUN /usr/local/bin/tinyclaw onboard

# --- tinyland customizations ---

# Bake config template with Aperture API routing placeholders.
# entrypoint.sh substitutes ANTHROPIC_API_KEY and ANTHROPIC_BASE_URL at startup.
COPY --chown=tinyclaw:tinyclaw tinyland/config.json /home/tinyclaw/.tinyclaw/config.json
COPY --chown=tinyclaw:tinyclaw tinyland/entrypoint.sh /usr/local/bin/entrypoint.sh

# Workspace bootstrap files — copied to /workspace-defaults/ so the K8s init
# container can seed the PVC on first boot without overwriting evolved state.
COPY --chown=tinyclaw:tinyclaw tinyland/workspace/ /workspace-defaults/

ENTRYPOINT ["entrypoint.sh"]
CMD ["gateway"]
