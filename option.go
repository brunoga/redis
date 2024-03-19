package redis

import "time"

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

type retryDelayOption time.Duration

func (opt retryDelayOption) apply(rwLock *RWLock) {
	rwLock.retryDelay = time.Duration(opt)
}

func WithRetryDelay(retryDelay time.Duration) Option {
	return retryDelayOption(retryDelay)
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
