package swapor

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

const (
	patchAeroPoolsKey     = "patchAeroPools"
	patchInfusionPoolsKey = "patchInfusionPools"
)

// patch aerodrome v2 pools stable & fee.
func (eh *EventHandler) PatchV2Factory(key string) {
	logger.Info().Msgf("PatchV2Factory: patch %s ....", key)

	var fac *common.SwapFactory
	if key == patchAeroPoolsKey {
		fac = &common.SwapFactory{
			Address: strings.ToLower("0x420DD381b31aEf6683db6B902084cB0FFECe40Da"),
			Name:    "AeroV2",
			Typ:     common.PoolTypeAeroAMM,
		}
	} else {
		fac = &common.SwapFactory{
			Address:   strings.ToLower("0x2d9a3a2bd6400ee28d770c7254ca840c82faf23f"),
			Name:      "Infusion",
			Typ:       common.PoolTypeInfusionAMM,
			StableFee: 500,
			Fee:       500,
		}
	}

	pairQuery := eh.pairsQuery

	pools, err := eh.pp.GetV2PoolInfos(pairQuery, fac, 1000, false)
	if err != nil {
		logger.Fatal().Err(err).
			Str("factory", fac.Address).
			Str("name", fac.Name).
			Msg("PatchAeroV2: get aero pool info failed")
	}

	err = eh.FlushToRedis(pools)
	if err != nil {
		logger.Fatal().Err(err).
			Str("factory", fac.Address).
			Str("name", fac.Name).
			Msg("PatchV2Factory: flush pools to redis failed")
	}

	logger.Info().
		Str("factory", fac.Address).
		Str("name", fac.Name).
		Int("pools", len(pools)).
		Msg("PatchV2Factory: load pool from chain success")
}

func (eh *EventHandler) FlushToRedis(pls map[string]*pool.Pool) error {
	pipe := eh.store.Pipeline()
	ctx := context.Background()

	for addr, pl := range pls {
		buf, err := json.Marshal(pl)
		if err != nil {
			logger.Error().Err(err).Msg("marshal pool failed")
			return err
		}

		if err := pipe.HSet(ctx, common.PoolLiquidKey, addr, buf).Err(); err != nil {
			logger.Error().Err(err).Str("pool", addr).Msg("HSET pool update failed")
			return err
		}
		// remove pool from unknownPoolInfo
		pipe.HDel(ctx, unknownPoolInfoKey, addr)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("sync pools updates to redis failed")
		return err
	}

	return nil
}

func (eh *EventHandler) DoPatches() {
	eh.doPatch(patchAeroPoolsKey)
	eh.doPatch(patchInfusionPoolsKey)
}

func (eh *EventHandler) doPatch(key string) {
	ctx := context.Background()

	v, err := eh.store.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			switch key {
			case patchAeroPoolsKey:
				fallthrough
			case patchInfusionPoolsKey:
				eh.PatchV2Factory(key)
				eh.store.Set(ctx, key, time.Now().Format(time.RFC3339), 0)
				logger.Info().Msgf("do patch %s successfully", key)

			default:
				logger.Error().Msgf("doPatch: unknown patch key: %s", key)
			}
		} else {
			logger.Fatal().Err(err).Str("key", key).Msg("doPatch: get redis key failed")
		}
	} else {
		logger.Info().Msgf("patch %s already done at: %v", key, v)
	}
}
