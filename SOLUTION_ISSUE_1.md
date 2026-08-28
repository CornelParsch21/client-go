# Solution for Issue #1

## 🛠️ Proposed Solution (by Aditya Waghamare)

### Analysis
`NewTokenBucketRateLimiter` could be initialized with a non‑positive burst, which makes the underlying `rate.NewLimiter` block forever because a token bucket of size 0 never yields tokens. This leads to indefinite request throttling.

### Fix
Sanitize the `burst` argument: if `qps > 0` and `burst <= 0`, coerce `burst` to `1` and emit a warning via `klog`. The change is isolated to the constructor and does not affect valid configurations.

### Implementation
```go
// util/flowcontrol/throttle.go
package flowcontrol

import (
    "golang.org/x/time/rate"
    "k8s.io/klog/v2"
)

type RateLimiter interface {
    Accept()
    TryAccept() bool
    // other methods …
}

// NewTokenBucketRateLimiter creates a token‑bucket based RateLimiter.
// It guarantees a minimum burst of 1 when qps > 0 to avoid indefinite blocking.
func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
    if qps > 0 && burst <= 0 {
        klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
        burst = 1
    }
    // Existing behaviour – keep the original limiter implementation.
    limiter := rate.NewLimiter(rate.Limit(qps), burst)
    return &tokenBucketRateLimiter{limiter: limiter}
}

// tokenBucketRateLimiter wraps golang.org/x/time/rate.Limiter to satisfy the
// local RateLimiter interface.
type tokenBucketRateLimiter struct {
    limiter *rate.Limiter
}

func (t *tokenBucketRateLimiter) Accept() {
    _ = t.limiter.Wait(context.Background())
}

func (t *tokenBucketRateLimiter) TryAccept() bool {
    return t.limiter.Allow()
}

// … any other methods that already existed remain unchanged.
```

### Testing
Add two unit tests to verify the safeguard works and does not alter normal behaviour.
```go
// util/flowcontrol/throttle_test.go
package flowcontrol

import "testing"

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
    rl := NewTokenBucketRateLimiter(10, 0) // burst will be coerced to 1
    if !rl.TryAccept() {
        t.Fatalf("expected TryAccept to succeed after burst coercion")
    }
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
    rl := NewTokenBucketRateLimiter(5, -5) // burst will be coerced to 1
    if !rl.TryAccept() {
        t.Fatalf("expected TryAccept to succeed with negative burst coerced to 1")
    }
}
```
Run `go test ./...` – both tests should pass and no request should block indefinitely.

**Signed‑off‑by:** Aditya Waghamare <adityawaghamare7620@gmail.com>

---
*Submitted by Aditya Waghamare*
💰 **Payout Address (Base L2 / EVM):** `0xb61dBcdBc3407F71EaCb64D4CBFAcf9FFfe2415C`