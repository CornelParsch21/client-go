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
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"k8s.io/klog/v2"
)

type RateLimiter interface {
	// TryAccept returns true if a token is taken immediately. Otherwise,
	// it returns false.
	TryAccept() bool
	// Accept returns once a token becomes available.
	Accept()
	// Stop stops the rate limiter, cleaning up any state associated with it.
	Stop()
	// QPS returns QPS of this rate limiter
	QPS() float32
}

type tokenBucketRateLimiter struct {
	limiter *rate.Limiter
	clock   Clock
	qps     float32
}

// NewTokenBucketRateLimiter creates a rate limiter which implements RateLimiter.
func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
	if qps > 0 && burst <= 0 {
		klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
		burst = 1
	}
	return &tokenBucketRateLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
		clock:   RealClock{},
		qps:     qps,
	}
}

type alwaysSafeRateLimiter struct{}

func (alwaysSafeRateLimiter) TryAccept() bool { return true }
func (alwaysSafeRateLimiter) Accept()         {}
func (alwaysSafeRateLimiter) Stop()           {}
func (alwaysSafeRateLimiter) QPS() float32    { return 0 }

// NewAlwaysSafeRateLimiter creates a rate limiter which always permits operations.
func NewAlwaysSafeRateLimiter() RateLimiter {
	return alwaysSafeRateLimiter{}
}

func (t *tokenBucketRateLimiter) TryAccept() bool {
	return t.limiter.AllowN(t.clock.Now(), 1)
}

// Accept will block until a token becomes available
func (t *tokenBucketRateLimiter) Accept() {
	t.limiter.Wait(context.Background())
}

func (t *tokenBucketRateLimiter) Stop() {}

func (t *tokenBucketRateLimiter) QPS() float32 {
	return t.qps
}
