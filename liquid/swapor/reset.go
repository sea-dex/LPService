package swapor

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	starcomm "starbase.ag/liquidity/liquid/common"
	"starbase.ag/liquidity/liquid/pool"
	"starbase.ag/liquidity/pkg/logger"
)

func (eh *EventHandler) getMode() string {
	var (
		err error
		val string
		ctx = context.Background()
	)

	val, err = eh.store.Get(ctx, modeKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Warn().Msg("not found mode in redis, set to resume, run as reset mode")

			return "reset"
		} else {
			logger.Error().Err(err).Msg("load mode failed")
		}
	}

	if val != "reset" && val != "resume" {
		logger.Fatal().Msgf("invalid run mode: %s", val)
	}

	logger.Info().Msgf("---------- run as %s mode ----------", val)

	return val
}

func (eh *EventHandler) DoReset(ctx context.Context) {
	logger.Info().Msg("reset mode, do reset ....")

	// eh.reseting = true
	_ = eh.cleanRedis(ctx)

	tokenKnown := map[string]bool{}
	for addr := range eh.tokens {
		tokenKnown[addr] = true
	}

	eh.LoadFactoryPoolsAndTokens(ctx, tokenKnown)
	// eh.reseting = false

	logger.Info().Msg("reset mode complete")
}

func (eh *EventHandler) loadTokenInfo(ctx context.Context) error {
	data, err := eh.store.HGetAll(ctx, starcomm.TokenInfoKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		logger.Error().Err(err).Str("key", starcomm.TokenInfoKey).Msg("Redis HGETALL failed")

		return err
	}

	for addr, val := range data {
		var tok starcomm.Token
		if err := json.Unmarshal([]byte(val), &tok); err != nil {
			logger.Warn().Err(err).Str("token", addr).Str("val", val).Msg("unmarshal token failed")
			continue
		}

		eh.tokens[strings.ToLower(addr)] = &tok
	}

	logger.Info().Msgf("load tokens from redis: %d", len(eh.tokens))

	return nil
}

func (eh *EventHandler) cleanRedis(ctx context.Context) error {
	logger.Info().Msg("reset mode, clean redis ....")

	// we don't clean tokenInfo
	if err := eh.store.Del(ctx, unknownFactory).Err(); err != nil {
		return err
	}

	if err := eh.store.Del(ctx, unknownPoolInfoKey).Err(); err != nil {
		return err
	}

	if err := eh.store.Del(ctx, starcomm.PoolLiquidKey).Err(); err != nil {
		return err
	}

	if err := eh.store.Del(ctx, eventsParsedAtKey).Err(); err != nil {
		return err
	}

	return nil
}

// LoadFactoryPools load factory pools/pairs, tokens.
func (eh *EventHandler) LoadFactoryPoolsAndTokens(ctx context.Context, tokenKnown map[string]bool) {
	tokens := map[string]bool{}

	t1 := time.Now().UnixMilli()

	logger.Info().Msg("start load factory pools ....")

	for _, factory := range eh.factory {
		pools, err := eh.loadFactoryPools(factory)
		if err != nil {
			for _, pl := range pools {
				if !tokenKnown[pl.Token0] {
					tokens[pl.Token0] = true
				}

				if !tokenKnown[pl.Token1] {
					tokens[pl.Token1] = true
				}
			}

			for _, p := range pools {
				eh.pools[p.Address] = p
			}
		}
	}

	t2 := time.Now().UnixMilli()
	logger.Info().Msgf("Load all factory pools success, used: %d ms", t2-t1)

	// send tokens to main thread
	logger.Info().Msgf("start load %d token info....", len(tokens))

	t3 := time.Now().UnixMilli()
	logger.Info().Msgf("load %d token info used: %d ms", len(tokens), t3-t2)
}

func (eh *EventHandler) loadFactoryPools(factory *starcomm.SwapFactory) (map[string]*pool.Pool, error) {
	switch factory.Typ {
	case starcomm.PoolTypeAMM:
		fallthrough
	case starcomm.PoolTypeInfusionAMM:
		fallthrough
	case starcomm.PoolTypeAeroAMM:
		return eh.getV2FactoryPools(factory)

	case starcomm.PoolTypeCAMM:
		fallthrough
	case starcomm.PoolTypePancakeCAMM:
		fallthrough
	case starcomm.PoolTypeAeroCAMM:
		return eh.getV3FactoryPools(factory)

	default:
		logger.Fatal().Msgf("invalid factory type: %d", factory.Typ)
	}

	panic("no reach here")
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (eh *EventHandler) getV3PoolsByEvents(factory []string, from, to uint64) (map[string][]string, error) {
	filterAddr := []common.Address{}

	for i, addr := range factory {
		factory[i] = strings.ToLower(addr)
		filterAddr = append(filterAddr, common.HexToAddress(addr))
	}

	sub := eh.subscriber
	if to == 0 {
		to = sub.GetLatestBlockNumber()
	}

	evts, err := sub.GetLogsFromToParallel(context.Background(), filterAddr, from, to, 5000)
	if err != nil {
		return nil, err
	}

	pools := map[string][]string{}

	for _, evt := range evts {
		if len(evt.Topics) == 0 {
			continue
		}

		topic := strings.ToLower(evt.Topics[0].String())
		if topic == pool.TopicPoolCreated {
			evtAddr := strings.ToLower(evt.Address.String())

			var poolAddr common.Address

			poolAddr.SetBytes(evt.Data[len(evt.Data)-20:])
			pools[evtAddr] = append(pools[evtAddr], poolAddr.String())
		}
	}

	return pools, nil
}

func (eh *EventHandler) getV2FactoryPools(factory *starcomm.SwapFactory) (map[string]*pool.Pool, error) {
	return eh.pp.GetV2PoolInfos(eh.pairsQuery, factory, 0, true)
}

func (eh *EventHandler) getV3FactoryPools(factory *starcomm.SwapFactory) (map[string]*pool.Pool, error) {
	var (
		err   error
		pools []string
	)

	t1 := time.Now().UnixMilli()

	if factory.Typ == starcomm.PoolTypeAeroCAMM {
		pools, err = eh.pp.GetAeroV3Pools(eh.pairsQuery, factory.Address, 0, true)
	} else {
		if factory.PositionManager == "" {
			logger.Fatal().
				Str("factory", factory.Address).
				Str("vender", factory.Name).
				Msg("factory positionManager is empty")
		}

		pools, err = eh.pp.GetV3Pools(eh.pairsQuery, factory.Address, factory.PositionManager, factory.Name, 0, true)
	}

	if err != nil {
		logger.Error().Err(err).Str("name", factory.Name).Msg("get factory pools failed")
		return nil, err
	}

	t2 := time.Now().UnixMilli()
	logger.Info().Str("factory", factory.Name).Msgf("getV3Pools: pools: %d, used %d ms", len(pools), t2-t1)

	wg := &sync.WaitGroup{}
	lk := &sync.Mutex{}
	failed := []string{}
	poolsMap := map[string]*pool.Pool{}
	routines := 0
	maxRoutines := 100
	maxRetry := 5
	getPoolFn := func(addr string) {
		var (
			pl  *pool.Pool
			err error
		)

		for retry := 0; retry < maxRetry; retry++ {
			pl, err = eh.getPool(addr)
			if err == nil {
				break
			}
		}

		lk.Lock()
		if err == nil {
			poolsMap[addr] = pl
		} else {
			failed = append(failed, addr)
		}
		lk.Unlock()
		wg.Done()
	}

	for _, addr := range pools {
		addr = strings.ToLower(addr)

		wg.Add(1)

		go getPoolFn(addr)

		routines++
		if routines >= maxRoutines {
			wg.Wait()

			routines = 0
		}
	}

	wg.Wait()

	t3 := time.Now().UnixMilli()
	logger.Info().Str("factory", factory.Name).Msgf("load all v3 pools liquidity used: %d ms, failed: %d", t3-t2, len(failed))

	return poolsMap, nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func (eh *EventHandler) resetPoolLiquidity() {
	ts1 := time.Now().UnixMilli()
	wg := &sync.WaitGroup{}
	failed := []string{}
	ctx := context.Background()
	counter := 0

	for _, pl := range eh.pools {
		wg.Add(1)
		pl.ResetLiquidity()

		go func(p *pool.Pool) {
			if err := eh.getPoolLiquidity(p); err != nil {
				// delete the pool temporary
				logger.Warn().
					Err(err).
					Str("pool", p.Address).
					Str("poolType", p.Typ.String()).
					Msg("reset mode: init pool liquidity from chain failed")

				// logger.Fatal().Msg("get pool liquidity failed: pool=%s error=%v", p.Address, err)
				failed = append(failed, p.Address)
			} else {
				buf, err := json.Marshal(p)
				if err != nil {
					logger.Error().Err(err).Str("pool", p.Address).Msg("resetPoolLiquidity: marshal pool failed")
				} else {
					if err := eh.store.HSet(ctx, starcomm.PoolLiquidKey, p.Address, buf).Err(); err != nil {
						logger.Error().Err(err).Str("pool", p.Address).Msg("HSET pool init liquidity failed")
					}
				}
			}

			wg.Done()
		}(pl)

		counter++

		if counter >= 100 {
			wg.Wait()
			logger.Info().Msgf("get pool liquidity done, count=%d", counter)
			counter = 0
		}
	}

	wg.Wait()

	// todo send all pool liquidity to kafka

	ts2 := time.Now().UnixMilli()
	logger.Info().Msgf("initialize pools liquidity onchain used %d ms, pools: %d failed: %d",
		ts2-ts1, len(eh.pools), len(failed))

	for _, item := range failed {
		delete(eh.pools, item)
		logger.Info().Str("pool", item).Msg("remove pool because of get liquidity failed")
	}
}
