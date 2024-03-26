package retrier

import "time"

// Fixed implements a fixed delay Retrier.
type Fixed struct {
	delay time.Duration
}

var _ Retrier = (*Fixed)(nil)

// NewFixed creates a new FixedRetrier instance.
func NewFixed(delay time.Duration) *Fixed {
	return &Fixed{
		delay: delay,
	}
}

// NextDelay returns a fixed delay.
func (f *Fixed) NextDelay() time.Duration {
	return f.delay
}
