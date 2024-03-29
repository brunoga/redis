package redis

import (
	"time"

	"github.com/brunoga/redis/retrier"
)

type Option interface {
	apply(*RWLock)
}

type keyTTLOption time.Duration

func (opt keyTTLOption) apply(rwLock *RWLock) {
	rwLock.keyTTL = time.Duration(opt)
}

func WithKeyTTL(keyTTL time.Duration) Option {
	return keyTTLOption(keyTTL)
}

type maxAttemptsOption uint8

func (opt maxAttemptsOption) apply(rwLock *RWLock) {
	rwLock.maxAttempts = uint8(opt)
}

func WithMaxAttempts(maxAttempts uint8) Option {
	return maxAttemptsOption(maxAttempts)
}

type autoRefreshOption bool

func (opt autoRefreshOption) apply(rwLock *RWLock) {
	rwLock.autoRefresh = bool(opt)
}

func WithAutoRefresh(autoRefresh bool) Option {
	return autoRefreshOption(autoRefresh)
}

type retrierOption struct {
	retrier retrier.Retrier
}

func (opt retrierOption) apply(rwLock *RWLock) {
	rwLock.retrier = opt.retrier
}

func WithRetrier(retrier retrier.Retrier) Option {
	return retrierOption{
		retrier: retrier,
	}
}
