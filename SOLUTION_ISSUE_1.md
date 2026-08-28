# Solution for Issue #1

## 🛠️ Proposed Solution (by Aditya Waghamare)

### Analysis
In `client-go`'s `util/flowcontrol` package, `tokenBucketRateLimiter` wraps `golang.org/x/time/rate.Limiter`. When a rate limiter is configured with `qps > 0` but `burst <= 0`, `rate.NewLimiter` creates a token bucket that can never hold tokens, causing `Wait()` / `Accept()` calls to block indefinitely and `TryAccept()` to permanently return `false`.

To prevent indefinite request throttling and potential client deadlocks caused by invalid configuration, `NewTokenBucketRateLimiter` must sanitize the `burst` argument when `qps > 0`, coercing any value `<= 0` to a minimum safe burst size of `1` while logging a warning via `klog`.

---

### Implementation

#### `util/flowcontrol/throttle.go`
```go
package flowcontrol

import (
	"golang.org/x/time/rate"
	"k8s.io/klog/v2"
)

// NewTokenBucketRateLimiter creates a RateLimiter implementation based on token bucket algorithm.
func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
	if qps > 0 && burst <= 0 {
		klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
		burst = 1
	}
	return &tokenBucketRateLimiter{
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
)

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, 0)
	if !limiter.TryAccept() {
		t.Errorf("expected TryAccept() to return true when burst size 0 is coerced to 1")
	}
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, -5)
	if !limiter.TryAccept() {
		t.Errorf("expected TryAccept() to return true when negative burst size is coerced to 1")
	}
}
```

---

### Testing & Verification
1. Executed `go test ./util/flowcontrol/...` to confirm that `TestTokenBucketRateLimiter_ZeroBurst` and `TestTokenBucketRateLimiter_NegativeBurst` pass without blocking.
2. Verified that existing token bucket tests with valid QPS and burst configurations (e.g. infinite QPS or positive burst values) continue to pass without regression.
3. Confirmed `klog` correctly logs warnings when `burst <= 0` with `qps > 0`.

Signed-off-by: Aditya Waghamare <adityawaghamare7620@gmail.com>

---
*Submitted by Aditya Waghamare*
💰 **Payout Address (Base L2 / EVM):** `0xb61dBcdBc3407F71EaCb64D4CBFAcf9FFfe2415C`