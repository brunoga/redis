package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRWLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "localhost:6379",
	})

	rwLock := NewRWLock(client, "test", 1*time.Second)

	ctx := context.Background()

	if err := rwLock.Lock(ctx); err != nil {
		t.Fatalf("Lock failed: %v", err)
	}
	fmt.Println("Writer lock acquired.")

	go func() {
		time.Sleep(500 * time.Millisecond)
		if err := rwLock.Unlock(ctx); err != nil {
			t.Errorf("Unlock failed: %v", err)
		}
		fmt.Println("Writer lock unlocked.")
	}()

	if err := rwLock.RLock(ctx); err != nil {
		t.Errorf("Lock failed: %v", err)
	}
	fmt.Println("Reader lock acquired.")
	defer func() {
		err := rwLock.RUnlock(ctx)
		if err != nil {
			t.Fatalf("Unlock failed: %v", err)
		}
		fmt.Println("Reader lock unlocked.")
	}()
}
