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
	r := NewTokenBucketRateLimiter(1, 1)
	r.Accept() // consume the only token
	start := time.Now()
	r.Accept() // block for 1s
	end := time.Now()
	if end.Sub(start) < 900*time.Millisecond {
		t.Errorf("expected to block for 1s, blocked for %v", end.Sub(start))
	}
}

func TestPreventInfiniteThrottleWithSmallBurst(t *testing.T) {
	// QPS = 1, Burst = 0
	limiter := NewTokenBucketRateLimiter(1.0, 0)

	// Accept should not block indefinitely. It should resolve quickly.
	done := make(chan struct{})
	go func() {
		limiter.Accept()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("RateLimiter.Accept() blocked indefinitely with burst = 0")
	}
}
