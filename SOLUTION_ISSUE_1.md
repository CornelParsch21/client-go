# Solution for Issue #1

## 🛠️ Proposed Solution (by Aditya Waghamare)

### Analysis
In `client-go`, `NewTokenBucketRateLimiter` delegates rate-limiting to `golang.org/x/time/rate.NewLimiter(rate.Limit(qps), burst)`. When `burst <= 0` while `qps > 0`, the underlying token bucket capacity is zero, causing methods like `Accept()`, `TryAccept()`, or `Wait()` to block indefinitely or fail continuously as tokens can never be acquired.

### Fix
Sanitize the `burst` argument inside `NewTokenBucketRateLimiter` before instantiating the rate limiter. If `qps > 0` and `burst <= 0`, log a warning via `klog.Warningf` and coerce `burst` to `1`. Unit tests are added to verify zero and negative burst handling.

### Implementation

#### `util/flowcontrol/throttle.go`

```go
package flowcontrol

import (
	"golang.org/x/time/rate"
	"k8s.io/klog/v2"
)

// NewTokenBucketRateLimiter creates a rate limiter with given QPS and burst values.
// If QPS > 0 and burst <= 0, burst is coerced to 1 with a warning logged.
func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
	if qps > 0 && burst <= 0 {
		klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
		burst = 1
	}
	return &tokenBucketLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
		qps:     qps,
		burst:   burst,
	}
}
```

#### `util/flowcontrol/throttle_test.go`

```go
package flowcontrol

import (
	"testing"
	"time"
)

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10.0, 0)
	if limiter.QPS() != 10.0 {
		t.Errorf("Expected QPS 10.0, got %f", limiter.QPS())
	}
	// Verify TryAccept succeeds and does not block indefinitely
	if !limiter.TryAccept() {
		t.Errorf("Expected TryAccept() to return true when burst is coerced to 1")
	}
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10.0, -5)
	if limiter.QPS() != 10.0 {
		t.Errorf("Expected QPS 10.0, got %f", limiter.QPS())
	}
	// Verify Accept does not block indefinitely
	done := make(chan struct{})
	go func() {
		limiter.Accept()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatalf("Accept() blocked indefinitely for negative burst configuration")
	}
}
```

### Testing
1. Run unit tests in `util/flowcontrol`:
   ```bash
   go test -v ./util/flowcontrol -run "TestTokenBucketRateLimiter_.*"
   ```
2. Verify `TestTokenBucketRateLimiter_ZeroBurst` and `TestTokenBucketRateLimiter_NegativeBurst` pass and `klog` outputs the warning for non-positive burst values.
3. Verify existing rate-limiting tests remain unaffected.

Signed-off-by: Aditya Waghamare <adityawaghamare7620@gmail.com>

---
*Submitted by Aditya Waghamare*
💰 **Payout Address (Base L2 / EVM):** `0xb61dBcdBc3407F71EaCb64D4CBFAcf9FFfe2415C`