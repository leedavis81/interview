# Rate limiter — coding exercise

Implement a **sliding-window rate limiter**: given a maximum number of requests and a time window, decide whether each incoming request from a client should be allowed or denied. Only requests that fall inside the rolling window count toward the limit (not a fixed clock boundary).

Starter code is in `start/` for each language; a reference implementation is in `solution/`.

| Language | Start | Solution |
|----------|-------|----------|
| Go | `go/start/main.go` | `go/solution/main.go` |
| Python | `python/start/rate_limiter.py` | `python/solution/rate_limiter.py` |
| Java | `java/start/RateLimiter.java` | `java/solution/RateLimiter.java` |
| Node.js | `nodejs/start/rateLimiter.js` | `nodejs/solution/rateLimiter.js` |
| Rust | `rust/start/` | `rust/solution/` |

---

## Hard requirements

- **Per-client tracking** — each client has an independent request count.
- **Sliding time window** — only requests within the configured window count toward the limit.
- **Allow / deny** — return whether a request is allowed and record it when allowed.

---

## Additional requirements

- **Handle concurrent requests** (thread-safe).
- **Be memory-efficient** (prune expired timestamps).

---

## Follow-up

After the coding exercise, use [system-design.md](system-design.md) for the distributed rate limiter discussion.
