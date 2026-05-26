import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class RateLimiter {

    private final int maxRequests;
    private final long timeWindowMs;
    private final Map<String, List<Long>> clients = new HashMap<>();

    public RateLimiter(int maxRequests, long timeWindowMs) {
        this.maxRequests = maxRequests;
        this.timeWindowMs = timeWindowMs;
    }

    public synchronized boolean isRequestAllowed(String clientId) {
        long now = System.currentTimeMillis();
        long windowStart = now - timeWindowMs;

        List<Long> timestamps = clients.getOrDefault(clientId, new ArrayList<>());

        // Drop expired timestamps — only requests inside the sliding window count.
        // Entries are appended in order, so we keep those still >= windowStart.
        List<Long> inWindow = new ArrayList<>();
        for (long ts : timestamps) {
            if (ts >= windowStart) {
                inWindow.add(ts);
            }
        }

        // Limit decision: deny if the client already has maxRequests in the window.
        if (inWindow.size() >= maxRequests) {
            clients.put(clientId, inWindow);
            return false;
        }

        // Allowed — record this request and persist the pruned list.
        inWindow.add(now);
        clients.put(clientId, inWindow);
        return true;
    }

    public static void main(String[] args) {
        RateLimiter rl = new RateLimiter(2, 10_000);
        System.out.println(rl.isRequestAllowed("user_123")); // true  — 1st request
        System.out.println(rl.isRequestAllowed("user_123")); // true  — 2nd request
        System.out.println(rl.isRequestAllowed("user_123")); // false — limit reached
        System.out.println(rl.isRequestAllowed("user_456")); // true  — different client
    }
}
