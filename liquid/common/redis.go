package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	PoolLiquidKey = "poolLiquidity"
	TokenInfoKey  = "tokenInfo"
)

// CreatePoolRedisStore create redis store.
func CreatePoolRedisStore(isProd bool, addr string, password string, db int) (redis.Cmdable, error) {
	if isProd {
		clusterCli := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(addr, ","),
			Password: password,
		})

		_, err := clusterCli.Ping(context.Background()).Result()
		if err != nil {
			return nil, fmt.Errorf("redis ping: %w", err)
		}

		return clusterCli, nil
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})

	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return cli, nil
}
