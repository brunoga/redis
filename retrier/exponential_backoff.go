package retrier

import (
	"math/rand"
	"time"
)

// ExponentialBackoff implements an exponential backoff Retrier.
type ExponentialBackoff struct {
	minDelay   time.Duration
	maxDelay   time.Duration
	retryCount int
}

var _ Retrier = (*ExponentialBackoff)(nil)

// NewExponentialBackoff creates a new ExponentialBackoffRetrier instance.
func NewExponentialBackoff(minDelay,
	maxDelay time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		minDelay: minDelay,
		maxDelay: maxDelay,
	}
}

// NextDelay calculates the next delay using a exponential backoff with jitter.
func (e *ExponentialBackoff) NextDelay() time.Duration {
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
func (e *ExponentialBackoff) Reset() {
	e.retryCount = 0
}
