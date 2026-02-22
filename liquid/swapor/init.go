package swapor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

// Init do initialize
// 1. load pool info/liquidity from redis
// 2. load blockNumber, logIndex.
func (eh *EventHandler) Init(mode string, dryrun bool) (err error) {
	ts1 := time.Now().UnixMilli()
	ctx := context.Background()
	_ = eh.loadTokenInfo(ctx)

	// load blocknumber, logIndex where last parsed at
	mode = eh.loadLastEventsParsedAt(ctx, mode)
	if mode == "reset" && !dryrun {
		// eh.resetCh = make(chan any)
		// eh.DoReset(ctx)
		_ = eh.cleanRedis(ctx)

		if err = eh.store.Set(ctx, modeKey, "resume", 0).Err(); err != nil {
			logger.Error().Err(err).Msg("set mode to resume failed")
		}

		return nil
	}

	if !dryrun {
		eh.DoPatches()
	}

	if err = eh.loadRedisPools(common.PoolLiquidKey, false); err != nil {
		return
	}

	if err = eh.loadRedisPools(unknownPoolInfoKey, true); err != nil {
		return
	}

	if !dryrun {
		eh.removePoolFromUnknown()
	}

	// load unknown factory
	data, err := eh.store.HGetAll(ctx, unknownFactory).Result()
	if err != nil {
		if err == redis.Nil {
			return
		}

		if !dryrun {
			for k := range data {
				if _, ok := eh.factory[k]; ok {
					// this factory is known, delete it from redis
					eh.store.HDel(ctx, unknownFactory, k)
				} else {
					eh.unknownFactory[k] = true
				}
			}
		}
	}

	if !dryrun {
		eh.checkPoolTokens()
		eh.notifyReboot()
	}

	ts2 := time.Now().UnixMilli()
	logger.Info().Msgf("load pools from redis used %d ms, pools: %d", ts2-ts1, len(eh.pools))

	eh.logFactoryPools()

	eh.lpservice.SetTokenPools(eh.tokens, eh.pools)
	eh.lpservice.InitPools()
	eh.lpservice.UpdateArbPairs(eh.cfg.Arb.CalcAllTokens)

	return
}

// send poolType 9999 to downstream.
func (eh *EventHandler) notifyReboot() {
	// remove kafka dependency
	// buf, _ := json.Marshal(pool.PoolReloadAll)
	// eh.producer.ProduceMessage(eh.cfg.Kafka.Topic, [][]byte{buf})
}

// AddFactory add factory.
func (eh *EventHandler) AddFactory(addr, pm string, name string, typ common.PoolType, fee, stableFee int, known bool) {
	addr = strings.ToLower(addr)
	if !common.IsValidPoolType(typ) {
		logger.Fatal().Msgf("factory is NOT valid: factory=%s address=%s poolType=%v", name, addr, typ)
	}

	if stableFee == 0 {
		stableFee = 5
	}

	factory := &common.SwapFactory{
		Address:         addr,
		PositionManager: pm,
		Name:            name,
		Typ:             typ,
		Fee:             fee,
		StableFee:       stableFee,
		Known:           known,
	}

	eh.factory[addr] = factory

	logger.Info().
		Str("factoryName", name).
		Str("factoryAddr", addr).
		Uint("poolType", uint(typ)).
		Msg("add swap factory")
}

func (eh *EventHandler) logFactoryPools() {
	for _, fac := range eh.factory {
		logger.Info().Msgf("%s pools: %d", fac.Name, fac.Pools)
	}
}

func (eh *EventHandler) removePoolFromUnknown() {
	toRemoved := []string{}

	for addr := range eh.unknownPools {
		if _, ok := eh.pools[addr]; ok {
			logger.Info().Msgf("remove pool %s from unknown pools", addr)
			toRemoved = append(toRemoved, addr)
		}
	}

	ctx := context.Background()
	pipeline := eh.store.Pipeline()

	for _, addr := range toRemoved {
		pipeline.HDel(ctx, unknownPoolInfoKey, addr)
		delete(eh.unknownPools, addr)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("delete unknown pool in redis failed")
	} else {
		logger.Info().
			Msgf("delete %d unknown pool in redis", len(toRemoved))
	}
}

func (eh *EventHandler) loadLastEventsParsedAt(ctx context.Context, mode string) string {
	if mode == "reset" {
		provider := eh.pp.Get()

		b, err := provider.Client.BlockNumber(ctx)
		if err != nil {
			logger.Fatal().Msg("loadLastEventsParsedAt: get block number failed: " + err.Error())
		}
		// start from current block
		eh.blockNumber = b.Uint64()

		logger.Info().Msgf("reset mode, use blocknumber %d to start subscribe events", eh.blockNumber)

		return mode
	}

	var (
		blockNumber uint64
		logIdx      uint
	)

	v, err := eh.store.Get(ctx, eventsParsedAtKey).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Fatal().Msg("loadLastEventsParsedAt: get eventsParsedAtKey failed: " + err.Error())
		}
		// use the block number as last events parsed
		provider := eh.pp.Get()

		b, err := provider.Client.BlockNumber(ctx)
		if err != nil {
			logger.Fatal().Msg("get block number failed: " + err.Error())
		}
		// start from current block
		blockNumber = b.Uint64()
		mode = "reset"

		logger.Warn().Msgf("loadLastEventsParsedAt: No events has been parsed, switch to reset mode, use blocknumber %d", blockNumber)
	} else {
		vv := strings.Split(v, ",")
		if blockNumber, err = strconv.ParseUint(vv[0], 10, 64); err != nil {
			logger.Fatal().Msg(fmt.Sprintf("invalid redis blockNumber: %v %v", v, err))
		}

		var index uint64

		index, err = strconv.ParseUint(vv[1], 10, 64)
		if err != nil {
			logger.Fatal().Msg(fmt.Sprintf("invalid redis logIndex: %v %v", v, err))
		}

		logIdx = uint(index)
	}

	eh.blockNumber = blockNumber
	eh.logIndex = logIdx

	return mode
}

func (eh *EventHandler) loadRedisPools(key string, unknown bool) error {
	ctx := context.Background()

	data, err := eh.store.HGetAll(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		logger.Error().Err(err).Str("key", key).Msg("Redis HGETALL failed")

		return err
	}

	for addr, val := range data {
		if unknown {
			var info pool.PoolBasicInfo
			if err := json.Unmarshal([]byte(val), &info); err != nil {
				logger.Error().Err(err).Str("val", val).Msg("Unmarshal pool basic info failed")
				return err
			}

			if fac, ok := eh.factory[info.Factory]; ok {
				//
				logger.Info().
					Str("pool", info.Address).
					Str("factory", info.Factory).
					Str("factoryName", fac.Name).
					Msg("change unknown pool to known")

				if err := eh.store.HDel(ctx, key, info.Address).Err(); err != nil {
					logger.Warn().Err(err).Msgf("delete pool from unknown redis failed: %s", info.Address)
				}

				fac.Pools++
			} else {
				eh.unknownPools[addr] = &info
			}
		} else {
			var pl pool.Pool
			if err := json.Unmarshal([]byte(val), &pl); err != nil {
				logger.Error().Err(err).Str("val", val).Msg("Unmarshal pool failed")
				return err
			}

			fac, ok := eh.factory[pl.Factory]
			if !ok {
				logger.Warn().
					Str("pool", pl.Address).
					Str("factory", pl.Factory).
					Str("Vendor", pl.Vendor).
					Msg("change pool to unknown")
				// eh.store.HDel(ctx, key, pool.Address)
				eh.unknownPools[addr] = &pool.PoolBasicInfo{
					Address: addr,
					Token0:  pl.Token0,
					Token1:  pl.Token1,
					Factory: pl.Factory,
				}
			} else {
				fac.Pools++

				pl.Reload()
				eh.pools[addr] = &pl
			}
		}
	}

	logger.Info().Int("pools", len(data)).Str("redisKey", key).Msg("Load pool from redis success")

	return nil
}
