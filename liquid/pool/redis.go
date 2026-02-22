package pool

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/redis/go-redis/v9"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/pkg/logger"
)

// LoadAllPools load all pools from redis.
func LoadAllPools(ps redis.Cmdable) (map[string]*Pool, error) {
	m, err := ps.HGetAll(context.Background(), common.PoolLiquidKey).Result()
	if err != nil {
		return nil, err
	}

	pools := map[string]*Pool{}

	for k, v := range m {
		var liq Pool

		err := json.Unmarshal([]byte(v), &liq)
		if err != nil {
			logger.Error().Err(err).Msg("unmarshal pool liquidity failed")
			continue
		}

		liq.Reload()
		pools[strings.ToLower(k)] = &liq
	}

	return pools, nil
}

// LoadTokens load tokens from redis.
func LoadTokens(conn redis.Cmdable) (map[string]*common.Token, error) {
	data, err := conn.HGetAll(context.Background(), common.TokenInfoKey).Result()
	if err != nil {
		logger.Error().Err(err).Msg("Redis HGETALL TokenInfoKey failed")
		return nil, err
	}

	tokens := map[string]*common.Token{}

	for key, val := range data {
		var item common.Token
		if err := json.Unmarshal([]byte(val), &item); err != nil {
			logger.Error().Err(err).Str("key", key).Str("val", val).Msg("Unmarshal token failed")
			return nil, err
		}

		tokens[item.Address] = &item
	}

	return tokens, nil
}
