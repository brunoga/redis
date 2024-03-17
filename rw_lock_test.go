package redis

import (
	"context"
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

	for i := 0; i < 10; i++ {
		err := rwLock.RLock(context.Background())
		if err != nil {
			t.Error(err)
		}
		defer func() {
			err = rwLock.RUnlock(context.Background())
			if err != nil {
				t.Error(err)
			}
		}()
	}

	//err := rwLock.RUnlock(context.Background())
	//if err != nil {
	//	t.Error(err)
	//}
}
