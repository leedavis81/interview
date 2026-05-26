# Implement a sliding-window rate limiter that restricts how many requests a client
# can make within a rolling time window. For example, if the limit is 100 requests
# per hour, only requests in the last hour count toward the limit (not a fixed
# clock-hour boundary).
#
# Requirements:
# - Per-client tracking — each client has an independent request count.
# - Sliding time window — only requests within the configured window count toward the limit.
# - Allow / deny — return whether a request is allowed and record it when allowed.
# - Handle concurrent requests (thread-safe).
# - Be memory-efficient (prune expired timestamps).


class RateLimiter:
    # TODO: add __init__

    def is_request_allowed(self, client_id: str) -> bool:
        # TODO: implement sliding-window check
        return True


if __name__ == "__main__":
    pass
    # rl = ...
    # print(rl.is_request_allowed("user_123"))
    # print(rl.is_request_allowed("user_123"))
    # print(rl.is_request_allowed("user_123"))
    # print(rl.is_request_allowed("user_456"))
