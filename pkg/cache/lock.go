package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrLock       = errors.New("lock failed")
	ErrUnlock     = errors.New("unlock failed")
	ErrLockExists = errors.New("already locked")
)

// Lock TODO: Using a better redis lock library is better.
func (cli *RedisClient) Lock(ctx context.Context, key string, expiry time.Duration) error {
	ok, err := cli.SetNX(ctx, key, 1, expiry).Result()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLock, err)
	}

	if !ok {
		return ErrLockExists
	}

	return nil
}

func (cli *RedisClient) Unlock(ctx context.Context, key string) error {
	err := cli.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnlock, err)
	}

	return nil
}
