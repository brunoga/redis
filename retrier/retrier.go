package retrier

import "time"

// Retrier is the interface to be implemented by different retry algorithms.
type Retrier interface {
	// NextDelay returns the next delay to be used in a retry loop.
	NextDelay() time.Duration
}
