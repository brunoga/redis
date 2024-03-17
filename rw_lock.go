package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/brunoga/redis/internal"
	"github.com/redis/go-redis/v9"
)

const (
	readerCountPrefix = "reader_count_"
	writerFlagPrefix  = "writer_flag_"
)

// RWLock is a Redis-based implementation of a distributed read-write lock.
type RWLock struct {
	client     *redis.Client
	id         string
	expiration time.Duration
}

// NewRWLock creates a new RWLock instance.
func NewRWLock(client *redis.Client, id string,
	expiration time.Duration) *RWLock {
	return &RWLock{
		client:     client,
		id:         id,
		expiration: expiration,
	}
}

// Lock acquires a write lock.
func (r *RWLock) Lock(ctx context.Context) error {
	for {
		result, err := r.tryLock(ctx)
		if err != nil {
			return err
		}

		if result != -1 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

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

	return nil
}

// RLock acquires a read lock.
func (r *RWLock) RLock(ctx context.Context) error {
	for {
		result, err := r.tryRLock(ctx)
		if err != nil {
			return err
		}

		if result != -1 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

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

	return nil
}

func (r *RWLock) tryLock(ctx context.Context) (int, error) {
	return r.runRedisScript(ctx, internal.LockScript, []string{
		readerCountPrefix + r.id,
		writerFlagPrefix + r.id,
	}, r.expiration.Milliseconds())
}

func (r *RWLock) tryRLock(ctx context.Context) (int, error) {
	return r.runRedisScript(ctx, internal.RLockScript, []string{
		readerCountPrefix + r.id,
		writerFlagPrefix + r.id,
	}, r.expiration.Milliseconds())
}

func (r *RWLock) unlock(ctx context.Context) (int, error) {
	return r.runRedisScript(ctx, internal.UnlockScript, []string{
		writerFlagPrefix + r.id,
	})
}

func (r *RWLock) rUnlock(ctx context.Context) (int, error) {
	return r.runRedisScript(ctx, internal.RUnlockScript, []string{
		readerCountPrefix + r.id,
	})
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
