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
)

func TestReader(t *testing.T) {
	tb := NewTokenBucketRateLimiter(10, 10)
	r := NewMeasuredRateLimiter(tb)
	if r.QPS() != 10 {
		t.Errorf("Expected QPS to be 10, got %f", r.QPS())
	}
}

func TestTokenBucketRateLimiter_ZeroBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, 0)
	if !limiter.TryAccept() {
		t.Error("Expected TryAccept to return true for coerced burst size of 1")
	}
}

func TestTokenBucketRateLimiter_NegativeBurst(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, -5)
	if !limiter.TryAccept() {
		t.Error("Expected TryAccept to return true for coerced burst size of 1")
	}
}
