package util

import (
    "golang.org/x/time/rate"
)

// NewRateLimiter creates a rate limiter with a minimum burst of 1.
// This prevents indefinite throttling when burst is configured too low.
func NewRateLimiter(qps float64, burst int) *rate.Limiter {
    // Ensure minimum burst of 1 to prevent indefinite throttling
    if burst < 1 {
        burst = 1
    }
    return rate.NewLimiter(rate.Limit(qps), burst)
}

// ValidateBurst checks that burst is sufficient for expected concurrency.
// Returns a recommended burst value if the provided burst is too low.
func ValidateBurst(burst, expectedConcurrency int) int {
    if burst < 1 {
        return 1
    }
    if burst < expectedConcurrency {
        return expectedConcurrency
    }
    return burst
}
