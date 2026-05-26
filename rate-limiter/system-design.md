# Distributed rate limiter — system design

Use this after the candidate has implemented the in-process sliding-window rate limiter.

---

## Scenario

You’ve just built a single-process rate limiter. Now imagine we deploy it behind a load balancer with **five application servers**. Requests are distributed in **round-robin**, so the same client’s traffic can land on different servers over time:

```
                         ┌──────────────┐      ┌─────────────────────┐
              Clients ──►│ Load         │      │ App Server 1        │
                         │ Balancer     │─────►│ App Server 2        │
                         │              │      │ App Server 3        │
                         └──────────────┘      │ App Server 4        │
                               round-robin      │ App Server 5        │
                                                └─────────────────────┘
                                                  application servers
```

> **Say to the candidate:** “With this setup, the in-process limiter on each server only sees part of a client’s traffic, so per-client limits won’t hold globally. How would you design a **distributed rate limiter** so each client still gets the **same limit across the whole fleet**, no matter which server handles the request? Walk me through your approach.”

---

## Key questions

### 1. Scaling the central store

**Ask:** The shared datastore is becoming a bottleneck. How would you scale it, and what would the architecture look like at higher load?

**Listen for:**

- **Partitioning / sharding** — e.g. by `clientID` (or hash of it) so each shard owns a slice of clients; consistent hashing if shards are added or removed.
- **Routing** — how an app server knows which shard to talk to (client-id hash, lookup service, proxy layer).
- **Hot keys** — awareness that a few very active clients can still saturate a single shard; mitigations (local caching, rate-limit tiers, spreading hot clients).
- **Read/write patterns** — whether they separate reads from writes, use replicas for reads, or push more work to the edge (local token cache with periodic sync).
- Tradeoffs between **strong global accuracy** and **approximate limits** under extreme scale.

---

### 2. Failure and replication

**Ask:** One of the database nodes has failed. How do you protect the data, and how do you keep traffic and ownership balanced across the remaining nodes?

**Listen for:**

- **Replication factor** — e.g. RF=3 with quorum reads/writes; what happens when a node is down (degraded but available vs hard stop).
- **Failover** — automatic promotion of replicas, health checks, and brief unavailability vs split-brain risk.
- **Rebalancing** — reassigning shards when capacity changes; avoiding thundering herds during redistribution.
- **Durability vs latency** — sync vs async replication and impact on rate-limit correctness if a primary fails mid-write.
- Whether they mention **backup / restore** or only in-memory/ephemeral rate-limit state (TTL data may be acceptable to lose briefly).

---

### 3. Observability

**Ask:** What metrics and signals would you watch to know the distributed rate limiter is healthy and stable?

**Listen for:**

- **Limit effectiveness** — rate of 429s, allow/deny ratio, per-client or per-tier breakdowns.
- **Store health** — latency p50/p99, error rates, connection pool exhaustion, shard CPU/memory, hot-key alerts.
- **App tier** — time spent in rate-limit checks, timeouts talking to the store, cache hit rate if they use local caches.
- **Correctness proxies** — sudden spikes in allows (possible bypass) or denials (misconfiguration).
- **SLO-oriented** — dashboards/alerts tied to user impact, not only infrastructure CPU.

---

### 4. Abuse at window boundaries

**Ask:** A client has figured out when the sliding window “resets” and sends large bursts at the end and start of each window. How would you defend against that, and would you change the algorithm?

**Listen for:**

- **Sliding vs fixed window** — why edge bursting happens with fixed windows; sliding window reduces but does not eliminate burstiness at the boundary.
- **Leaky bucket / token bucket** — smoother output rate; tradeoff of burst allowance vs strictness.
- **Dual limits** — e.g. 100/hour plus 10/minute to cap short bursts.
- **Implementation impact** — moving from timestamp lists to token refill logic; atomic decrement in the store; different Redis structures.
- Honest tradeoffs: **fairness**, **UX** (stricter = more 429s), and **operational complexity**.

---

### 5. Client experience when limited

**Ask:** When a client is rate limited, what should the API return, and what information helps them back off correctly?

**Listen for:**

- **HTTP 429** with a clear error body (machine-readable code + human message).
- **`Retry-After`** — seconds until they can retry; whether it’s wall-clock or aligned to window reset.
- **Rate-limit headers** — e.g. `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset` (or `RateLimit-*` standard headers); consistency across endpoints.
- **Idempotency / safety** — clients should not hammer retries; exponential backoff vs honoring `Retry-After`.
- Whether limits are **per API key, IP, user, or tenant** and that this is documented for integrators.

---

## Steering questions

Use these to guide the discussion if they get stuck or stay too high-level.

- Where does the request count / timestamp data live now that it can't be in-process memory?
- How would you structure the keys and data? (Look for sensible key schemes like `ratelimit:{clientID}` with sorted sets or counters with TTLs.)
- If two requests from the same client hit two different servers at the exact same moment, what happens?
- How do you handle stale data?
