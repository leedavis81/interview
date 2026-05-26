// Implement a sliding-window rate limiter that restricts how many requests a client
// can make within a rolling time window.
//
// Requirements:
// - Per-client tracking — each client has an independent request count.
// - Sliding time window — only requests within the configured window count toward the limit.
// - Allow / deny — return whether a request is allowed and record it when allowed.
// - Handle concurrent requests (thread-safe).
// - Be memory-efficient (prune expired timestamps).

pub struct RateLimiter {
    // TODO: fields
}

impl RateLimiter {
    // TODO: add a constructor — decide how to configure max requests and time window
    // pub fn new(...) -> Self {
    //     Self { }
    // }

    pub fn is_request_allowed(&mut self, _client_id: &str) -> bool {
        // TODO: implement sliding-window check
        true
    }
}

fn main() {
    // let mut rl = RateLimiter::new(/* your configuration */);
    // println!("{}", rl.is_request_allowed("user_123"));
    // println!("{}", rl.is_request_allowed("user_123"));
    // println!("{}", rl.is_request_allowed("user_123"));
    // println!("{}", rl.is_request_allowed("user_456"));
}
