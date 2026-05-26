class RateLimiter {
  /**
   * @param {number} maxRequests
   * @param {number} timeWindowMs
   */
  constructor(maxRequests, timeWindowMs) {
    this.maxRequests = maxRequests;
    this.timeWindowMs = timeWindowMs;
    /** @type {Map<string, number[]>} */
    this.clients = new Map();
  }

  /**
   * @param {string} clientId
   * @returns {boolean}
   */
  isRequestAllowed(clientId) {
    const now = Date.now();
    const windowStart = now - this.timeWindowMs;

    let timestamps = this.clients.get(clientId) ?? [];

    // Drop expired timestamps — only requests inside the sliding window count.
    // Entries are appended in order, so we keep those still >= windowStart.
    timestamps = timestamps.filter((ts) => ts >= windowStart);

    // Limit decision: deny if the client already has maxRequests in the window.
    if (timestamps.length >= this.maxRequests) {
      this.clients.set(clientId, timestamps);
      return false;
    }

    // Allowed — record this request and persist the pruned list.
    timestamps.push(now);
    this.clients.set(clientId, timestamps);
    return true;
  }
}

const rl = new RateLimiter(2, 10_000);
console.log(rl.isRequestAllowed("user_123")); // true  — 1st request
console.log(rl.isRequestAllowed("user_123")); // true  — 2nd request
console.log(rl.isRequestAllowed("user_123")); // false — limit reached
console.log(rl.isRequestAllowed("user_456")); // true  — different client

module.exports = { RateLimiter };
