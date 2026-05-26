# Sliding-window rate limiter — reference solution.

import threading
import time


class RateLimiter:
    def __init__(self, max_requests: int, time_window_seconds: float):
        self.max_requests = max_requests
        self.time_window_seconds = time_window_seconds
        self._lock = threading.Lock()
        self._clients: dict[str, list[float]] = {}

    def is_request_allowed(self, client_id: str) -> bool:
        with self._lock:
            now = time.time()
            window_start = now - self.time_window_seconds

            timestamps = self._clients.get(client_id, [])

            # Drop expired timestamps — only requests inside the sliding window count.
            # Entries are appended in order, so we keep those still >= window_start.
            timestamps = [ts for ts in timestamps if ts >= window_start]

            # Limit decision: deny if the client already has max_requests in the window.
            if len(timestamps) >= self.max_requests:
                self._clients[client_id] = timestamps
                return False

            # Allowed — record this request and persist the pruned list.
            timestamps.append(now)
            self._clients[client_id] = timestamps
            return True


if __name__ == "__main__":
    rl = RateLimiter(max_requests=2, time_window_seconds=10)
    print(rl.is_request_allowed("user_123"))  # True  — 1st request
    print(rl.is_request_allowed("user_123"))  # True  — 2nd request
    print(rl.is_request_allowed("user_123"))  # False — limit reached
    print(rl.is_request_allowed("user_456"))  # True  — different client
