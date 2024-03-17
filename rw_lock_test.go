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

	go func() {
		for {
			simulateWorkAsync(rwLock, 100*time.Millisecond, false, 5)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for {
			simulateWorkAsync(rwLock, 100*time.Millisecond, true, 1)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(10 * time.Second)
}

func simulateWorkAsync(lock *RWLock, sleep time.Duration, write bool, count int) {
	for i := 0; i < count; i++ {
		go func() {
			ctx := context.Background()
			if write {
				if err := lock.Lock(ctx); err != nil {
					fmt.Printf("Lock failed: %v\n", err)
				}
				fmt.Println("Writer lock acquired.")
				time.Sleep(sleep)
				if err := lock.Unlock(ctx); err != nil {
					fmt.Printf("Unlock failed: %v\n", err)
				}
				fmt.Println("Writer lock unlocked.")
			} else {
				if err := lock.RLock(ctx); err != nil {
					fmt.Printf("Lock failed: %v\n", err)
				}
				fmt.Println("Reader lock acquired.")
				time.Sleep(sleep)
				if err := lock.RUnlock(ctx); err != nil {
					fmt.Printf("Unlock failed: %v\n", err)
				}
				fmt.Println("Reader lock unlocked.")
			}
		}()
	}
}
