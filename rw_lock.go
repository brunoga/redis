package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brunoga/redis/internal"
	"github.com/redis/go-redis/v9"
)

const (
	readerCountKeyPrefix = "reader_count_"
	writerCountKeyPrefix = "writer_count_"
)

// RWLock is a Redis-based implementation of a distributed read-write lock.
type RWLock struct {
	client     redis.Scripter
	id         string
	expiration time.Duration

	readerCountKey string
	writerCountKey string

	refreshCh chan struct{}

	m          sync.Mutex
	refreshing bool
}

// NewRWLock creates a new RWLock instance.
func NewRWLock(client redis.Scripter, id string,
	expiration time.Duration) *RWLock {
	return &RWLock{
		client:         client,
		id:             id,
		expiration:     expiration,
		readerCountKey: readerCountKeyPrefix + id,
		writerCountKey: writerCountKeyPrefix + id,
		refreshCh:      make(chan struct{}),
	}
}

// Lock acquires a write lock.
func (r *RWLock) Lock(ctx context.Context) error {
	err := lockLoop(ctx, func(ctx context.Context) (bool, error) {
		result, err := r.tryLock(ctx)
		if err != nil {
			return false, err
		}

		return result != -1, nil
	})
	if err != nil {
		return err
	}

	r.startRefreshLoop(ctx)

	return nil
}

// TryLock tries to acquire a write lock.
func (r *RWLock) TryLock(ctx context.Context) (bool, error) {
	result, err := r.tryLock(ctx)
	if err != nil {
		return false, err
	}

	if result == -1 {
		return false, nil
	}

	return true, nil
}

// Unlock releases a write lock.
func (r *RWLock) Unlock(ctx context.Context) error {
	result, err := r.unlock(ctx)
	if err != nil {
		return err
	}

	if result == -1 {
		return fmt.Errorf("too many unlocks")
	}

	r.stopRefreshLoop()

	return nil
}

// Key returns the writer count key for this lock.
func (r *RWLock) Key() string {
	return r.writerCountKey
}

// RLock acquires a read lock.
func (r *RWLock) RLock(ctx context.Context) error {
	err := lockLoop(ctx, func(ctx context.Context) (bool, error) {
		result, err := r.tryRLock(ctx)
		if err != nil {
			return false, err
		}

		return result != -1, nil
	})
	if err != nil {
		return err
	}

	r.startRefreshLoop(ctx)

	return nil
}

// TryRLock tries to acquire a read lock.
func (r *RWLock) TryRLock(ctx context.Context) (bool, error) {
	result, err := r.tryRLock(ctx)
	if err != nil {
		return false, err
	}

	if result == -1 {
		return false, nil
	}

	return true, nil
}

// RUnlock releases a read lock.
func (r *RWLock) RUnlock(ctx context.Context) error {
	result, err := r.rUnlock(ctx)
	if err != nil {
		return err
	}

	if result == -1 {
		return fmt.Errorf("too many unlocks")
	}

	r.stopRefreshLoop()

	return nil
}

// Key returns the writer count key for this lock.
func (r *RWLock) RKey() string {
	return r.readerCountKey
}

func (r *RWLock) tryLock(ctx context.Context) (int, error) {
	result, err := r.runRedisScript(ctx, internal.LockScript, []string{
		r.readerCountKey,
		r.writerCountKey,
	}, r.expiration.Milliseconds())

	if err == nil && result == -1 {
		// Refresh writer count.
		r.refresh(ctx, []string{r.writerCountKey, ""})
	}

	return result, err
}

func (r *RWLock) tryRLock(ctx context.Context) (int, error) {
	result, err := r.runRedisScript(ctx, internal.RLockScript, []string{
		r.readerCountKey,
		r.writerCountKey,
	}, r.expiration.Milliseconds())

	if err == nil && result == -1 {
		// Refresh reader count.
		r.refresh(ctx, []string{"", r.readerCountKey})
	}

	return result, err
}

func (r *RWLock) unlock(ctx context.Context) (int, error) {
	return r.runRedisScript(ctx, internal.UnlockScript, []string{
		r.writerCountKey,
	})
}

func (r *RWLock) rUnlock(ctx context.Context) (int, error) {
	return r.runRedisScript(ctx, internal.RUnlockScript, []string{
		r.readerCountKey,
	})
}

func (r *RWLock) refresh(ctx context.Context, keys []string) {
	r.runRedisScript(ctx, internal.RefreshScript, keys,
		r.expiration.Milliseconds())
}

func lockLoop(ctx context.Context, f func(context.Context) (bool, error)) error {
	// TODO(bga): Make this configurable.
	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ok, err := f(ctx)
			if err != nil {
				return err
			}

			if ok {
				return nil
			}
		}
	}
}

func (r *RWLock) startRefreshLoop(ctx context.Context) {
	r.m.Lock()
	if !r.refreshing {
		r.refreshing = true
		go r.refreshLoop(ctx)
	}
	r.m.Unlock()
}

func (r *RWLock) stopRefreshLoop() {
	r.m.Lock()
	if r.refreshing {
		r.refreshing = false
		r.refreshCh <- struct{}{}
	}
	r.m.Unlock()
}

func (r *RWLock) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(r.expiration / 2)
	defer ticker.Stop()

L:
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Refresh reader and writer counts.
			r.refresh(ctx, []string{
				r.readerCountKey,
				r.writerCountKey,
			})
		case <-r.refreshCh:
			break L
		}
	}
}

func (r *RWLock) runRedisScript(ctx context.Context, script *redis.Script,
	keys []string, args ...interface{}) (int, error) {
	cmd := script.Run(ctx, r.client, keys, args...)
	if cmd.Err() != nil {
		return -1, cmd.Err()
	}

	result, err := cmd.Int()
	if err != nil {
		return -1, fmt.Errorf("failed reading result as int: %w", err)
	}

	return result, nil
}
