package flowcontrol

import (
	"context"

	"golang.org/x/time/rate"
	"k8s.io/klog/v2"
)

type RateLimiter interface {
	TryAccept() bool
	Accept()
	Stop()
	QPS() float32
}

type tokenBucketRateLimiter struct {
	limiter *rate.Limiter
	qps     float32
}

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

func (t *tokenBucketRateLimiter) TryAccept() bool {
	return t.limiter.Allow()
}

func (t *tokenBucketRateLimiter) Accept() {
	t.limiter.Wait(context.Background())
}

func (t *tokenBucketRateLimiter) Stop() {}

func (t *tokenBucketRateLimiter) QPS() float32 {
	return t.qps
}
