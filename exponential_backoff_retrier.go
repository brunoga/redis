package redis

import (
	"math/rand"
	"time"
)

// ExponentialBackoffRetrier implements an exponential backoff Retrier.
type ExponentialBackoffRetrier struct {
	minDelay   time.Duration
	maxDelay   time.Duration
	retryCount int
}

var _ Retrier = (*ExponentialBackoffRetrier)(nil)

// NewExponentialBackoffRetrier creates a new ExponentialBackoffRetrier instance.
func NewExponentialBackoffRetrier(minDelay,
	maxDelay time.Duration) *ExponentialBackoffRetrier {
	return &ExponentialBackoffRetrier{
		minDelay: minDelay,
		maxDelay: maxDelay,
	}
}

// NextDelay calculates the next delay using a exponential backoff with jitter.
func (e *ExponentialBackoffRetrier) NextDelay() time.Duration {
	// Calculate exponential delay.
	expDelay := e.minDelay * time.Duration(1<<e.retryCount)

	// Add jitter to avoid the thundering herd problem.
	jitter := time.Duration(rand.Int63n(int64(expDelay))) / 2

	totalDelay := expDelay + jitter

	// Make sure we don't exceed the maximum delay
	if totalDelay > e.maxDelay {
		totalDelay = e.maxDelay
	}

	e.retryCount++

	return totalDelay
}

// Reset resets the retry count.
func (e *ExponentialBackoffRetrier) Reset() {
	e.retryCount = 0
}
