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
	// TODO(bga): Loop until lock is acquired or context is done.
	return r.TryLock(ctx)
}

// TryLock tries to acquire a write lock.
func (r *RWLock) TryLock(ctx context.Context) error {
	result, err := r.runRedisScript(ctx, internal.LockScript, []string{
		readerCountPrefix + r.id,
		writerFlagPrefix + r.id,
	}, r.expiration.Milliseconds())
	if err != nil {
		return err
	}

	if result == -1 {
		return fmt.Errorf("can´t acquire write lock")
	}

	return nil
}

// Unlock releases a write lock.
func (r *RWLock) Unlock(ctx context.Context) error {
	result, err := r.runRedisScript(ctx, internal.UnlockScript, []string{
		writerFlagPrefix + r.id,
	})
	if err != nil {
		return err
	}

	if result == -1 {
		return fmt.Errorf("can't release write lock")
	}

	return nil
}

// RLock acquires a read lock.
func (r *RWLock) RLock(ctx context.Context) error {
	// TODO(bga): Loop until lock is acquired or context is done.
	return r.TryRLock(ctx)
}

// TryRLock tries to acquire a read lock.
func (r *RWLock) TryRLock(ctx context.Context) error {
	result, err := r.runRedisScript(ctx, internal.RLockScript, []string{
		readerCountPrefix + r.id,
		writerFlagPrefix + r.id,
	}, r.expiration.Milliseconds())
	if err != nil {
		return err
	}

	if result == -1 {
		return fmt.Errorf("can´t acquire read lock")
	}

	fmt.Printf("Readers after rlock: %d\n", result)

	return nil
}

// RUnlock releases a read lock.
func (r *RWLock) RUnlock(ctx context.Context) error {
	result, err := r.runRedisScript(ctx, internal.RUnlockScript, []string{
		readerCountPrefix + r.id,
	})
	if err != nil {
		return err
	}

	if result == -1 {
		return fmt.Errorf("too many unlocks")
	}

	fmt.Printf("Readers after runlock: %d\n", result)

	return nil
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
