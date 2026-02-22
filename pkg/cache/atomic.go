package cache

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrEvalLuaScript             = errors.New("failed to eval Lua script")
	ErrLuaScriptUnexpectedResult = errors.New("unexpected result from Lua script")
)

func (cli *RedisClient) CheckAndSetKey(ctx context.Context, key, value string) (bool, error) {
	luaScript := `
		local key = KEYS[1]
		local value = ARGV[1]

		local exists = redis.call('EXISTS', key)

		if exists == 0 then
			redis.call('SET', key, value)
			return 1
		else
			return 0
		end
	`

	result, err := cli.Eval(ctx, luaScript, []string{key}, value).Result()
	if err != nil {
		return false, fmt.Errorf("cache.CheckAndSetKey: %w, %w", ErrEvalLuaScript, err)
	}

	if intValue, ok := result.(int64); ok {
		return intValue == 1, nil
	}

	return false, fmt.Errorf("cache.CheckAndSetKey: %w, %w", ErrLuaScriptUnexpectedResult, err)
}
