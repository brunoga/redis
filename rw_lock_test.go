package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRWLock_RLock_RUnlock_Success(t *testing.T) {
	s := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	rwLock := NewRWLock(client, "test", 100*time.Millisecond)

	acquireReadAndCheckKeyValue(t, rwLock, s, "1")
	acquireReadAndCheckKeyValue(t, rwLock, s, "2")
	releaseReadAndCheckKeyValue(t, rwLock, s, "1")
	acquireReadAndCheckKeyValue(t, rwLock, s, "2")
	releaseReadAndCheckKeyValue(t, rwLock, s, "1")
	releaseReadAndCheckKeyValue(t, rwLock, s, "0")
}

func TestRWLock_RUnlock_NotLocked(t *testing.T) {
	s := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	rwLock := NewRWLock(client, "test", 100*time.Millisecond)

	err := rwLock.RUnlock(context.Background())
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	checkKeyValue(t, s, rwLock.RKey(), "0")
}

func TestRWLock_Lock_Unlock_Success(t *testing.T) {
	s := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	rwLock := NewRWLock(client, "test", 100*time.Millisecond)

	acquireWriteAndCheckKeyValue(t, rwLock, s, "1")
	releaseWriteAndCheckKeyValue(t, rwLock, s, "0")
}

func TestRWLock_Unlock_NotLocked(t *testing.T) {
	s := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    s.Addr(),
	})

	rwLock := NewRWLock(client, "test", 100*time.Millisecond)

	err := rwLock.Unlock(context.Background())
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	checkKeyValue(t, s, rwLock.Key(), "0")
}

func acquireReadAndCheckKeyValue(t *testing.T, rwLock *RWLock,
	s *miniredis.Miniredis, want string) {
	t.Helper()

	err := rwLock.RLock(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	checkKeyValue(t, s, rwLock.RKey(), want)
}

func releaseReadAndCheckKeyValue(t *testing.T, rwLock *RWLock,
	s *miniredis.Miniredis, want string) {
	t.Helper()

	err := rwLock.RUnlock(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	checkKeyValue(t, s, rwLock.RKey(), want)
}

func acquireWriteAndCheckKeyValue(t *testing.T, rwLock *RWLock,
	s *miniredis.Miniredis, want string) {
	t.Helper()

	err := rwLock.Lock(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	checkKeyValue(t, s, rwLock.Key(), want)
}

func releaseWriteAndCheckKeyValue(t *testing.T, rwLock *RWLock,
	s *miniredis.Miniredis, want string) {
	t.Helper()

	err := rwLock.Unlock(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	checkKeyValue(t, s, rwLock.Key(), want)
}

func checkKeyValue(t *testing.T, s *miniredis.Miniredis, key, want string) {
	t.Helper()

	got, err := s.Get(key)
	if err != nil {
		if want == "0" && err == miniredis.ErrKeyNotFound {
			got = "0"
		} else {
			t.Fatal(err)
		}
	}
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
