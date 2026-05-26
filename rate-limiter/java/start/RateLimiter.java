// Implement a sliding-window rate limiter that restricts how many requests a client
// can make within a rolling time window.
//
// Requirements:
// - Per-client tracking — each client has an independent request count.
// - Sliding time window — only requests within the configured window count toward the limit.
// - Allow / deny — return whether a request is allowed and record it when allowed.
// - Handle concurrent requests (thread-safe).
// - Be memory-efficient (prune expired timestamps).

public class RateLimiter {

    // TODO: add a constructor — decide how to configure max requests and time window

    public boolean isRequestAllowed(String clientId) {
        // TODO: implement sliding-window check
        return true;
    }

    public static void main(String[] args) {
        // RateLimiter rl = new RateLimiter(/* your configuration */);
        // System.out.println(rl.isRequestAllowed("user_123"));
        // System.out.println(rl.isRequestAllowed("user_123"));
        // System.out.println(rl.isRequestAllowed("user_123"));
        // System.out.println(rl.isRequestAllowed("user_456"));
    }
}
