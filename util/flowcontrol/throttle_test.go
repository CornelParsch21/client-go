/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flowcontrol

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	RateLimiter := NewTokenBucketRateLimiter(10, 1)
	RateLimiter.Accept() // consume the one token
	start := time.Now()
	RateLimiter.Accept() // block for 100ms
	end := time.Now()
	if end.Sub(start) < 90*time.Millisecond {
		t.Errorf("expected to block for 100ms, only blocked for %v", end.Sub(start))
	}
}

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, 0)
	if !limiter.TryAccept() {
		t.Error("expected TryAccept to return true for coerced burst of 1")
	}

	done := make(chan struct{})
	go func() {
		limiter.Accept()
		close(done)
	}()
	select {
	case <-done:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Accept blocked indefinitely with zero burst")
	}
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, -5)
	if !limiter.TryAccept() {
		t.Error("expected TryAccept to return true for coerced burst of 1")
	}

	done := make(chan struct{})
	go func() {
		limiter.Accept()
		close(done)
	}()
	select {
	case <-done:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Accept blocked indefinitely with negative burst")
	}
}
