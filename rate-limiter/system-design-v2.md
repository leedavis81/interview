# Distributed rate limiter — system design (v2)

Incremental interview flow. Use this after the candidate has implemented the in-process sliding-window rate limiter.

Each phase ends with a **constraint** that motivates the next. Don't jump ahead — if they mention sharding in phase 2, note it ("we'll come back to that") and keep them on the happy path until the current phase is solid.

---

## Interview arc

```
Phase 1  Recognize the gap          →  why in-process fails
Phase 2  Central store              →  shared state, data model
Phase 3  Correctness                →  races, atomicity, fail-open/closed
Phase 4  Latency & scale            →  edge cache, sharding, failure
Phase 5  Algorithm & client API     →  burst abuse, 429 / headers
Phase 6  Observability              →  metrics, alerts, correctness proxies
```

**Suggested timing**

| Track | Phases | ~Duration |
|-------|--------|-----------|
| Core path (mid-level) | 1 → 2 → 3 → 5b (429 only) | 25–30 min |
| Full path (senior) | 1 → 6 | 40–50 min |
| Stretch | Hot keys, dual limits, per-tenant tiers | as time allows |

---

## Setup

### Prerequisite

Candidate has built a working **in-process sliding-window** rate limiter (e.g. timestamp list or equivalent).

### Opening scenario

Five application servers behind a load balancer. Requests are distributed **round-robin**, so the same client's traffic lands on different servers over time:

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

> **Say to the candidate:** "You've deployed your rate limiter on each of these servers. The limit is 100 requests per minute per client. Walk me through what happens when a single client sends 100 requests in one minute — and then what happens if they keep sending."

---

## Phase 1 — Recognize the gap

**Goal:** Candidate feels why per-server state breaks global limits before you ask for a distributed design.

### Ask

1. With five servers and round-robin, does each client still get **at most 100 requests per minute** globally? Why or why not?
2. If not, what's the worst case? Can you put a number on it?

### Listen for

- Each server maintains its **own** counter; a client can get up to **~5× the limit** (500/min with a 100/min cap and 5 servers).
- Round-robin makes the split roughly even over time, but any distribution model that spreads one client across servers has the same problem.
- Recognition that the fix requires **shared state** or **centralized coordination** — not "sticky sessions" alone (sticky helps but doesn't solve fleet restarts, scaling, or uneven stickiness).

### Bridge to phase 2

> "OK — we need one logical limit across the fleet. Where does the request count live now that it can't be in-process memory? Walk me through one request end-to-end: what gets read, what gets written, and when."

**If stuck:** "What if every server called the same external service before allowing a request?"

---

## Phase 2 — Central store

**Goal:** Baseline distributed architecture — usually a single shared datastore (Redis is the common answer).

### Ask

1. What store would you use, and why?
2. How would you represent a **sliding window** for client `X` in that store?
3. Show me the keys and values after client `X` has made 3 requests in the last 30 seconds.

### Listen for

- **Redis** (or Memcached, DynamoDB, etc.) as a fast, external key-value store accessible from all app servers.
- **Key scheme** — e.g. `ratelimit:{clientID}`, optionally suffixed with window or algorithm type.
- **Data structures:**
  - Sorted set of timestamps (sliding window) — `ZADD` + `ZREMRANGEBYSCORE` + `ZCARD`
  - Counter + window start (fixed window) — simpler but burstier; note if they already implemented sliding in-process
  - Single counter with TTL — fixed window approximation
- **TTL** on keys so stale clients don't accumulate forever.
- **Identity** — what is `clientID`? API key, user ID, IP, tenant; mention this can be a separate discussion.

### Steering questions

- Where does the request count / timestamp data live?
- How would you structure the keys and data?

### Bridge to phase 3

> "Two requests from the same client hit two different servers at the **exact same moment**. Both servers read the count, both see 99/100, both allow. What happened?"

---

## Phase 3 — Correctness under concurrency

**Goal:** Atomic check-and-increment; explicit policy when the store is unavailable.

### Ask

1. How do you make the limit check **atomic** so concurrent requests can't race?
2. The shared store is down for 30 seconds. Do you allow traffic through or block it? Why?
3. How long are you willing to wait on the store before deciding?

### Listen for

- **Atomic operations** — Redis `INCR`, Lua script wrapping read-check-write, `MULTI`/`EXEC`, or compare-and-swap patterns.
- **Lua script example (conceptual):** read count, if under limit then increment and return allow, else return deny — all in one round-trip.
- **Fail-open vs fail-closed:**
  - Fail-open (allow when store is down) — protects availability; risk of abuse during outage.
  - Fail-closed (deny when store is down) — protects the system; risk of blocking legitimate traffic.
  - Reasonable answer depends on context (public API vs internal service); candidate should articulate the tradeoff, not just pick one.
- **Timeout budget** — e.g. 5–10ms for rate-limit check; on timeout, apply the same fail-open/closed policy consistently.
- Optional mention of **idempotency** — rate-limit check should be side-effect-safe on retry.

### Bridge to phase 4

> "Every allowed request now pays for a network round-trip to the store before your handler runs. Your API SLO is 100ms p99. The store adds 15ms p99. Is that acceptable? How would you reduce the cost?"

---

## Phase 4 — Latency, availability, and scale

**Goal:** Edge optimization first, then horizontal scale of the store, then failure handling. Three beats in one phase — take them in order.

### 4a — Reduce per-request cost

**Ask:** How do you keep rate limiting from dominating request latency?

**Listen for:**

- **Local token cache / edge allowance** — server pre-fetches or holds a small local budget, syncs with central store periodically.
- **Approximate limits** — intentionally looser global enforcement in exchange for fewer store calls; candidate should acknowledge the tradeoff.
- **Connection pooling**, pipelining, colocating Redis near app tier.
- **Async / non-blocking** — only if they can explain how the request still gets a synchronous allow/deny decision.

### 4b — Scale the central store

**Ask:** The shared datastore is becoming a bottleneck — CPU, memory, or network. How would you scale it, and what does the architecture look like at higher load?

**Listen for:**

- **Partitioning / sharding** — by `clientID` hash so each shard owns a slice of clients.
- **Routing** — app server hashes `clientID` → shard; consistent hashing if shards are added or removed.
- **Hot keys** — a few very active clients can saturate one shard; mitigations (local cache, dedicated tier, spreading hot clients).
- **Read/write patterns** — replicas for read scaling (with staleness caveats), or push more work to the edge.
- Tradeoff between **strong global accuracy** and **approximate limits** under extreme scale.

### 4c — Failure and replication

**Ask:** One of the database nodes has failed. How do you protect the data, and how do you keep traffic and ownership balanced across the remaining nodes?

**Listen for:**

- **Replication factor** — e.g. RF=3 with quorum reads/writes; degraded but available vs hard stop.
- **Failover** — automatic promotion, health checks, brief unavailability vs split-brain risk.
- **Rebalancing** — reassigning shards when capacity changes; avoiding thundering herds during redistribution.
- **Durability vs latency** — sync vs async replication; impact on correctness if a primary fails mid-write.
- Rate-limit state is often **ephemeral** (TTL data) — brief loss or reset may be acceptable; candidate should say so explicitly.

### Bridge to phase 5

> "The system is up, limits are mostly accurate, latency is acceptable. A client has figured out when the sliding window resets and sends huge bursts at the boundary. Is that a problem? Would you change the algorithm?"

---

## Phase 5 — Algorithm and client experience

**Goal:** Refine limit *behavior* and the *contract* with API consumers. Two beats — algorithm first, then HTTP surface.

### 5a — Abuse at window boundaries

**Ask:** A client sends large bursts at the end and start of each window. How would you defend against that, and would you change the algorithm?

**Listen for:**

- **Sliding vs fixed window** — fixed windows are especially vulnerable to 2× bursts at boundaries; sliding window reduces but doesn't eliminate burstiness.
- **Token bucket / leaky bucket** — smoother output rate; burst allowance as a tunable parameter.
- **Dual limits** — e.g. 100/hour **and** 10/minute to cap short bursts.
- **Implementation impact** — moving from timestamp lists to token refill logic; atomic decrement in the store; different Redis structures.
- Tradeoffs: **fairness**, **UX** (stricter = more 429s), **operational complexity**.

### 5b — Client experience when limited

**Ask:** When a client is rate limited, what should the API return, and what information helps them back off correctly?

**Listen for:**

- **HTTP 429** with a clear error body (machine-readable code + human message).
- **`Retry-After`** — seconds until retry; wall-clock vs aligned to window reset.
- **Rate-limit headers** — `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset` (or `RateLimit-*` standard headers); consistency across endpoints.
- **Backoff guidance** — honor `Retry-After` vs exponential backoff; don't hammer retries.
- Limits are **per API key, IP, user, or tenant** — documented for integrators.

### Bridge to phase 6

> "This is in production. How do you know the rate limiter is healthy, correct, and not silently failing open?"

---

## Phase 6 — Observability

**Goal:** Metrics and alerts tied to user impact, not just infrastructure.

### Ask

What metrics and signals would you watch to know the distributed rate limiter is healthy and stable? What would you alert on?

### Listen for

- **Limit effectiveness** — rate of 429s, allow/deny ratio, per-client or per-tier breakdowns.
- **Store health** — latency p50/p99, error rates, connection pool exhaustion, shard CPU/memory, hot-key alerts.
- **App tier** — time spent in rate-limit checks, timeouts talking to the store, cache hit rate if using local caches.
- **Correctness proxies** — sudden spikes in allows (possible bypass / fail-open) or denials (misconfiguration).
- **SLO-oriented** — dashboards and alerts tied to user impact, not only infrastructure CPU.

---

## Steering questions (any phase)

Use when the candidate is stuck or stays too high-level.

| Phase | Prompt |
|-------|--------|
| 1 | "If the limit is 100/min and there are 5 servers, what's the maximum a single client could send in one minute?" |
| 2 | "Draw me the Redis keys for client `abc` over a 60-second window." |
| 3 | "What happens if two servers both read `count=99` at the same time?" |
| 4 | "What if Redis latency goes from 2ms to 200ms?" |
| 4 | "What if one Redis shard owns 80% of traffic?" |
| 5 | "What's the difference between a fixed and sliding window at the `:00` boundary?" |
| 5 | "What should the client do when it gets a 429?" |
| 6 | "Allows doubled overnight but traffic didn't — what do you check first?" |

---

## Interviewer notes

### Keep one problem per phase

Resist letting them run ahead. Good sign: they name a later concern early. Response: "Good — hold that thought; we'll get there."

### Make them quantify

Numbers create inevitability. Examples:

- 5 servers × 100/min limit → up to **500/min** without shared state.
- Store timeout of 10ms on a 100ms SLO → **10%** of budget on rate limiting alone.

### Use their vocabulary

If they said "sorted set," the phase 3 question becomes "two `ZADD`s at once — what's the count?" If they said "token bucket," phase 5 builds on their refill logic.

### Repeat the core tradeoff

**Stronger global accuracy vs lower latency / higher availability** — revisit at phases 3, 4, and 5. Each time they should go deeper.

### Stretch topics (if ahead of schedule)

- Per-tenant tiers (free vs paid limits).
- Regional / multi-datacenter limits (eventual consistency across regions).
- Rate limiting vs DDoS protection (different layer, different tools).
- Cost of Redis memory at scale (millions of active keys).
- Sticky sessions as a partial mitigation — and why it's insufficient alone.

### Calibration

| Signal | Level |
|--------|-------|
| Phases 1–2 solid, fuzzy on atomicity | Mid |
| Phases 1–3 solid, thoughtful fail-open/closed | Senior |
| Phase 4 with hot-key awareness + phase 5 dual limits | Staff+ |
| Jumps to Kafka / service mesh without establishing shared state | Redirect to phase 2 |

---

## Quick reference — question chain

One-line prompts in interview order:

1. **Gap** — "Five servers, round-robin — does 100/min still hold?"
2. **Store** — "Where does the count live? Walk me through one request."
3. **Race** — "Same client, two servers, same millisecond — what happens?"
4. **Latency** — "Store adds 15ms p99. Acceptable? How do you shrink it?"
5. **Scale** — "Redis is the bottleneck. Now what?"
6. **Failure** — "A shard primary dies mid-write. What happens?"
7. **Abuse** — "Client games the window boundary. Change the algorithm?"
8. **Contract** — "They got a 429. What do they see in the response?"
9. **Ops** — "How do you know it's working in production?"
