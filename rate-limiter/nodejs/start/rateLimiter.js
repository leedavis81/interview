// Implement a sliding-window rate limiter that restricts how many requests a client
// can make within a rolling time window.
//
// Requirements:
// - Per-client tracking — each client has an independent request count.
// - Sliding time window — only requests within the configured window count toward the limit.
// - Allow / deny — return whether a request is allowed and record it when allowed.
// - Handle concurrent requests (thread-safe in Node's single-threaded model; consider
//   async callers if you add locking for a multi-threaded runtime).
// - Be memory-efficient (prune expired timestamps).

class RateLimiter {
  // TODO: add a constructor

  /**
   * @param {string} clientId
   * @returns {boolean}
   */
  isRequestAllowed(clientId) {
    // TODO: implement sliding-window check
    return true;
  }
}

// const rl = ...
// console.log(rl.isRequestAllowed("user_123"));
// console.log(rl.isRequestAllowed("user_123"));
// console.log(rl.isRequestAllowed("user_123"));
// console.log(rl.isRequestAllowed("user_456"));

module.exports = { RateLimiter };
