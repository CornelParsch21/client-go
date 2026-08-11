package flowcontrol

import (
	"testing"
)

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
	rl := NewTokenBucketRateLimiter(10, 0)
	if rl.QPS() != 10 {
		t.Errorf("expected QPS to be 10, got %f", rl.QPS())
	}

	if !rl.TryAccept() {
		t.Error("expected TryAccept to succeed with coerced burst")
	}
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
	rl := NewTokenBucketRateLimiter(10, -5)
	if rl.QPS() != 10 {
		t.Errorf("expected QPS to be 10, got %f", rl.QPS())
	}

	if !rl.TryAccept() {
		t.Error("expected TryAccept to succeed with coerced negative burst")
	}
}

func TestTokenBucketRateLimiter_ValidBurst(t *testing.T) {
	rl := NewTokenBucketRateLimiter(10, 5)
	if rl.QPS() != 10 {
		t.Errorf("expected QPS to be 10, got %f", rl.QPS())
	}

	for i := 0; i < 5; i++ {
		if !rl.TryAccept() {
		t.Errorf("expected TryAccept to succeed on iteration %d", i)
		}
	}
}
