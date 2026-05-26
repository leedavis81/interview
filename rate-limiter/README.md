# Rate limiter (coding + system design)

Per-language **start** (stub) and **solution** (reference) for a sliding-window rate limiter. Follow-up discussion uses [system-design.md](system-design.md).

## Layout

| Language | Start | Solution |
|----------|-------|----------|
| Go | `go/start/main.go` | `go/solution/main.go` |
| Python | `python/start/rate_limiter.py` | `python/solution/rate_limiter.py` |
| Java | `java/start/RateLimiter.java` | `java/solution/RateLimiter.java` |
| Node.js | `nodejs/start/rateLimiter.js` | `nodejs/solution/rateLimiter.js` |
| Rust | `rust/start/` | `rust/solution/` |

## Run reference solutions

```bash
# Go
cd go/solution && go run .

# Python
python3 python/solution/rate_limiter.py

# Java
cd java/solution && javac RateLimiter.java && java RateLimiter

# Node.js
node nodejs/solution/rateLimiter.js

# Rust
cd rust/solution && cargo run
```

Expected output (four lines): `true`, `true`, `false`, `true`.
