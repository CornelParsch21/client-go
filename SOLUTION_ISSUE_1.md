# Solution for Issue #1

## 🛠️ Proposed Solution (by Aditya Waghamare)

### Analysis
In `client-go`, `util/flowcontrol` wraps `golang.org/x/time/rate.NewLimiter`. When a `RateLimiter` is initialized with `burst <= 0` while `qps > 0`, the underlying token bucket has a capacity of zero or negative tokens. As a result, tokens can never be consumed, causing `Wait()` or `Accept()` to block indefinitely and `TryAccept()` to fail continuously.

### Fix
Sanitize the `burst` configuration parameter inside `NewTokenBucketRateLimiter` (and `NewTokenBucketRateLimiterWithClock`). When `qps > 0` and `burst <= 0`, emit a warning log using `klog.Warningf` and coerce `burst` to a minimum safe threshold of `1`.

### Implementation

#### File: `util/flowcontrol/throttle.go`

```go
package flowcontrol

import (
	"golang.org/x/time/rate"
	"k8s.io/klog/v2"
)

// NewTokenBucketRateLimiter creates a rate limiter that allows up to qps requests per second,
// with a maximum burst size of burst.
func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
	if qps > 0 && burst <= 0 {
		klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
		burst = 1
	}
	return &tokenBucketRateLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
		qps:     qps,
	}
}

// NewTokenBucketRateLimiterWithClock creates a rate limiter with a custom clock.
func NewTokenBucketRateLimiterWithClock(qps float32, burst int, c clock.Clock) RateLimiter {
	if qps > 0 && burst <= 0 {
		klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
		burst = 1
	}
	return &tokenBucketRateLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
		qps:     qps,
		clock:   c,
	}
}
```

#### File: `util/flowcontrol/throttle_test.go`

```go
package flowcontrol

import (
	"testing"
)

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
	qps := float32(10)
	burst := 0

	limiter := NewTokenBucketRateLimiter(qps, burst)
	if limiter == nil {
		t.Fatalf("expected non-nil RateLimiter")
	}

	if !limiter.TryAccept() {
		t.Errorf("expected TryAccept to return true after burst coercion")
	}
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
	qps := float32(10)
	burst := -5

	limiter := NewTokenBucketRateLimiter(qps, burst)
	if limiter == nil {
		t.Fatalf("expected non-nil RateLimiter")
	}

	if !limiter.TryAccept() {
		t.Errorf("expected TryAccept to return true after burst coercion")
	}
}
```

### Testing
1. Run `go test ./util/flowcontrol/...` to verify `TestTokenBucketRateLimiter_ZeroBurst` and `TestTokenBucketRateLimiter_NegativeBurst` pass without deadlocks.
2. Verify existing flowcontrol unit tests (`TestTokenBucketRateLimiter`, `TestPassiveRateLimiter`) continue to pass.

Signed-off-by: Aditya Waghamare <adityawaghamare7620@gmail.com>

---
*Submitted by Aditya Waghamare*
💰 **Payout Address (Base L2 / EVM):** `0xb61dBcdBc3407F71EaCb64D4CBFAcf9FFfe2415C`