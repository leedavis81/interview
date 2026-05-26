use std::collections::HashMap;
use std::sync::Mutex;
use std::time::{Duration, Instant};

pub struct RateLimiter {
    max_requests: usize,
    time_window: Duration,
    clients: Mutex<HashMap<String, Vec<Instant>>>,
}

impl RateLimiter {
    pub fn new(max_requests: usize, time_window: Duration) -> Self {
        Self {
            max_requests,
            time_window,
            clients: Mutex::new(HashMap::new()),
        }
    }

    pub fn is_request_allowed(&self, client_id: &str) -> bool {
        let mut clients = self.clients.lock().unwrap();
        let now = Instant::now();

        let timestamps = clients.entry(client_id.to_string()).or_default();

        // Drop expired timestamps — only requests inside the sliding window count.
        // Entries are appended in order, so retain keeps those still within time_window.
        timestamps.retain(|ts| now.duration_since(*ts) < self.time_window);

        // Limit decision: deny if the client already has max_requests in the window.
        if timestamps.len() >= self.max_requests {
            return false;
        }

        // Allowed — record this request (pruned list is already in the map).
        timestamps.push(now);
        true
    }
}

fn main() {
    let rl = RateLimiter::new(2, Duration::from_secs(10));
    println!("{}", rl.is_request_allowed("user_123")); // true  — 1st request
    println!("{}", rl.is_request_allowed("user_123")); // true  — 2nd request
    println!("{}", rl.is_request_allowed("user_123")); // false — limit reached
    println!("{}", rl.is_request_allowed("user_456")); // true  — different client
}
