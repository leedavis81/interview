// Implement a sliding window rate limiter which restricts how many requests a given client
// can make within a rolling time window. For example, if the limit is 100 requests per hour,
// a client who made their first request at 14:05 can make up to 100 requests in the
// window 14:05–15:05 (not just within a fixed clock-hour boundary)
//
// Requirements:
// - Per-client tracking - Each has its own independent request count.
// - Sliding time window — only requests made within the last time window duration count towards the limit.
// - Allow / deny decision function

package main

// RateLimiter enforces a sliding-window request limit per client.
type RateLimiter struct {
	// TODO: fields
}

// TODO: add a constructor — decide how to configure max requests and time window
// func NewRateLimiter(...) *RateLimiter {
// 	return &RateLimiter{}
// }

func (rl *RateLimiter) IsRequestAllowed(clientID string) bool {
	// TODO: implement sliding-window check
	return true
}

func main() {
	// rl := NewRateLimiter(/* your configuration */)
	// fmt.Println(rl.IsRequestAllowed("user_123"))
	// fmt.Println(rl.IsRequestAllowed("user_123"))
	// fmt.Println(rl.IsRequestAllowed("user_123"))
	// fmt.Println(rl.IsRequestAllowed("user_456"))
}
