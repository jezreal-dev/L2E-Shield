# L2E-Shield Go Proxy

L2E-Shield is a high-performance Go proxy gateway designed to shield LLM endpoints from malicious prompts and denial-of-wallet exhaustion attacks.

## Key Capabilities

1. Payload Validation: Intercepts incoming requests and evaluates prompt contents for injection patterns and jailbreak bypass vectors before calling downstream APIs.
2. Redis ZSET Rate Limiting: Restricts client requests using a Redis-backed sliding window filter to block token exhaustion.
3. SHA256 Response Caching: Computes hashes of incoming payloads and caches successful completions in Redis to prevent redundant charges.

## Directory Layout

1. cmd/server/main.go: Server bootstrap and router configurations.
2. internal/config/config.go: Environment variables schema.
3. internal/cache/redis.go: Redis integration client wrapper.
4. internal/security/detector.go: Security inspection checks.
5. internal/middleware/logging.go: Structural logger middleware.
6. internal/middleware/ratelimit.go: ZSET rate limiter middleware.
7. internal/handler/proxy.go: Main client proxy request handler.

## Local Setup Instructions

### Running with Docker Compose

1. Create a copy of the environmental variables template:
   cp .env.example .env
2. Edit .env and supply your GEMINI_API_KEY.
3. Launch the compose services:
   docker compose up --build

### Running Manually

1. Supply environment variables directly:
   export GEMINI_API_KEY="your_api_key"
   export REDIS_URL="redis://localhost:6379/0"
2. Compile and launch:
   go build -o server cmd/server/main.go
   ./server
