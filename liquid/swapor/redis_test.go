package swapor

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/utils"
)

func TestNotExistKey(t *testing.T) {
	utils.SkipCI(t)

	c, err := common.CreatePoolRedisStore(false, "127.0.0.1:6379", "", 0)
	assert.Nil(t, err)

	ctx := context.Background()
	key := "aaaa"
	assert.Nil(t, c.Del(ctx, key).Err())

	val, err := c.Get(ctx, key).Result()
	assert.Equal(t, redis.Nil, err)
	assert.Equal(t, "", val)
}

/*
func TestRedisPipe(t *testing.T) {
	utils.SkipCI(t)

	c, err := CreatePoolRedisStore(false, "127.0.0.1:6379", "", 0)
	assert.Nil(t, err)

	ctx := context.Background()
	c.Del(ctx, "a1")
	c.Del(ctx, "a2")
	c.Del(ctx, "a3")

	pipe := c.Pipeline()
	err = pipe.Set(ctx, "a1", "a1", 0).Err()
	assert.Nil(t, err)

	err = pipe.Set(ctx, "a2", "a2", 0).Err()
	assert.Nil(t, err)
	panic("interrupt")

	err = pipe.Set(ctx, "a3", "a3", 0).Err()
	assert.Nil(t, err)

	_, err = pipe.Exec(ctx)
	assert.Nil(t, err)
}
*/
