package redis

import "time"

// FixedRetrier implements a fixed delay Retrier.
type FixedRetrier struct {
	delay time.Duration
}

var _ Retrier = (*FixedRetrier)(nil)

// NewFixedRetrier creates a new FixedRetrier instance.
func NewFixedRetrier(delay time.Duration) *FixedRetrier {
	return &FixedRetrier{
		delay: delay,
	}
}

// NextDelay returns a fixed delay.
func (f *FixedRetrier) NextDelay() time.Duration {
	return f.delay
}
